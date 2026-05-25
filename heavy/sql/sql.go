package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/heavy/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNotFound = sql.ErrNoRows

type Config struct {
	Driver      string
	WaitingTime int
	Leader      ConnConfig
	Follower    ConnConfig
}

type ConnConfig struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
	SSL      bool
	Schema   string
	Options  ConnOptions
	MockDB   *sql.DB
}

type ConnOptions struct {
	MaxLifeTime time.Duration
	MaxIdle     int
	MaxOpen     int
}

type Interface interface {
	Leader() Command
	Follower() Command
	Stop()
	Shutdown(ctx context.Context) error
}

type sqlDB struct {
	endOnce    *sync.Once
	leader     Command
	follower   Command
	cfg        Config
	log        log.Interface
	shutdownMu sync.Mutex
	closing    bool
	inflight   sync.WaitGroup
}

func Init(cfg Config, log log.Interface) (Interface, error) {
	sql := &sqlDB{
		endOnce: &sync.Once{},
		log:     log,
		cfg:     cfg,
	}

	if err := sql.initDB(); err != nil {
		return nil, err
	}
	return sql, nil
}

func (s *sqlDB) Leader() Command {
	return s.leader
}

func (s *sqlDB) Follower() Command {
	return s.follower
}

func (s *sqlDB) Stop() {
	s.endOnce.Do(func() {
		ctx := context.Background()
		if s.leader != nil {
			if err := s.leader.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}
		if s.follower != nil {
			if err := s.follower.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}
	})
}

// Shutdown gracefully closes database connections, waiting for in-flight
// queries to complete until the context deadline is reached.
// If no deadline is set on the provided context, a default timeout of 30 seconds is used.
// Returns nil if all in-flight operations complete before the deadline.
// Returns a timeout error if the deadline is reached while operations are still active.
func (s *sqlDB) Shutdown(ctx context.Context) error {
	// Set closing flag to prevent new queries
	s.shutdownMu.Lock()
	s.closing = true
	s.shutdownMu.Unlock()

	// If no deadline is set, apply a default 30-second timeout
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Wait for in-flight operations to complete or context deadline to be reached
	done := make(chan struct{})
	go func() {
		s.inflight.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All in-flight operations completed, close connections cleanly
		s.Stop()
		return nil
	case <-ctx.Done():
		// Deadline reached, force-close connections
		s.Stop()
		return errors.NewWithCode(codes.CodeSQLShutdown, "shutdown timed out: %v", ctx.Err())
	}
}

func (s *sqlDB) initDB() error {
	ctx := context.Background()
	const maxRetries = 3

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(s.cfg.WaitingTime * int(time.Second)))
			s.log.Warn(ctx, fmt.Sprintf("SQL: [DB] retrying connection attempt %d/%d for db %s leader: %s on port %d", attempt, maxRetries, s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port))
		}

		db, err := s.connect(true)
		if err != nil {
			lastErr = err
			continue
		}

		s.log.Info(ctx, fmt.Sprintf("SQL: [LEADER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port, s.cfg.Leader.SSL))
		s.leader = initCommand(db, s.log)

		if s.isFollowerEnabled() {
			db, err = s.connect(false)
			if err != nil {
				lastErr = err
				continue
			}
			s.log.Info(ctx, fmt.Sprintf("SQL: [FOLLOWER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Follower.DB, s.cfg.Follower.Host, s.cfg.Follower.Port, s.cfg.Follower.SSL))
			s.follower = initCommand(db, s.log)
		} else {
			s.follower = s.leader
		}

		return nil
	}

	return errors.NewWithCode(codes.CodeSQLInit, "failed to connect to database after %d attempts: %v", maxRetries, lastErr)
}

func (s *sqlDB) connect(toLeader bool) (*sqlx.DB, error) {
	conf := s.cfg.Leader
	if !toLeader {
		conf = s.cfg.Follower
	}

	if toLeader {
		if s.cfg.Leader.MockDB != nil {
			return sqlx.NewDb(s.cfg.Leader.MockDB, s.cfg.Driver), nil
		}
	} else {
		if s.cfg.Follower.MockDB != nil {
			return sqlx.NewDb(s.cfg.Follower.MockDB, s.cfg.Driver), nil
		}
	}

	uri, err := s.getURI(conf)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(s.cfg.Driver, uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, errors.NewWithCode(codes.CodeSQLInit, err.Error())
	}

	sqlxDB := sqlx.NewDb(db, s.cfg.Driver)
	sqlxDB.SetMaxOpenConns(conf.Options.MaxOpen)
	sqlxDB.SetMaxIdleConns(conf.Options.MaxIdle)
	sqlxDB.SetConnMaxLifetime(conf.Options.MaxLifeTime)

	return sqlxDB, nil
}

func (s *sqlDB) isFollowerEnabled() bool {
	isHostNotEmpty := s.cfg.Follower.Host != ""
	isHostDifferent := (s.cfg.Follower.Host != s.cfg.Leader.Host && s.cfg.Follower.Port == s.cfg.Leader.Port)
	isPortDifferent := (s.cfg.Follower.Host == s.cfg.Leader.Host && s.cfg.Follower.Port != s.cfg.Leader.Port)
	return isHostNotEmpty && (isHostDifferent || isPortDifferent)
}

func (s *sqlDB) getURI(conf ConnConfig) (string, error) {
	switch s.cfg.Driver {
	case "postgres":
		ssl := `disable`
		if conf.SSL {
			ssl = `require`
		}
		if conf.Schema == "" {
			conf.Schema = "public"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.DB, conf.Schema, ssl), nil
	case "mysql":
		ssl := `false`
		if conf.SSL {
			ssl = `true`
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s&parseTime=true", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, ssl), nil
	default:
		return "", fmt.Errorf(`DB Driver [%s] is not supported`, s.cfg.Driver)
	}
}
