package sql

import (
	"context"
	"database/sql"

	"github.com/alpardfm/go-toolkit/log"
	"github.com/jmoiron/sqlx"
)

type Command interface {
	Close() error
	Ping(ctx context.Context) error
	In(query string, args ...interface{}) (string, []interface{}, error)
	Rebind(query string) string
	QueryIn(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error)
	Prepare(ctx context.Context, name string, query string) (CommandStmt, error)

	NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error)
	Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, name string, opts TxOptions) (CommandTx, error)

	Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error
}

type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

type command struct {
	db  *sqlx.DB
	log log.Interface
}

func initCommand(db *sqlx.DB, log log.Interface) Command {
	return &command{
		db:  db,
		log: log,
	}
}

func (c *command) Close() error {
	return c.db.Close()
}

func (c *command) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *command) In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}

func (c *command) Rebind(query string) string {
	return c.db.Rebind(query)
}

func (c *command) QueryIn(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	q, newArgs, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	return c.Query(ctx, name, c.Rebind(q), newArgs...)
}

// QueryRow should be avoided as it cannot be mocked using ExpectQuery
func (c *command) QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error) {
	row := c.db.QueryRowxContext(ctx, query, args...)
	return row, row.Err()
}

func (c *command) Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.db.QueryxContext(ctx, query, args...)
}

func (c *command) NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error) {
	return c.db.NamedQueryContext(ctx, query, arg)
}

func (c *command) Prepare(ctx context.Context, name string, query string) (CommandStmt, error) {
	stmt, err := c.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(ctx, name, stmt), nil
}

func (c *command) NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error) {
	return c.db.NamedExecContext(ctx, query, args)
}

func (c *command) Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *command) BeginTx(ctx context.Context, name string, opt TxOptions) (CommandTx, error) {
	opts := &sql.TxOptions{
		Isolation: opt.Isolation,
		ReadOnly:  opt.ReadOnly,
	}
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return initTx(ctx, name, tx, opts, c.log), nil
}

func (c *command) Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}
