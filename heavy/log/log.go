package log

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alpardfm/go-toolkit/appcontext"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/operator"
	"github.com/rs/zerolog"
)

var once = sync.Once{}

type Interface interface {
	// TODO add Debugf
	Trace(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	Info(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
	Fatal(ctx context.Context, obj any)
}

type Config struct {
	Level string
}

type logger struct {
	log zerolog.Logger
}

func Init(cfg Config) (Interface, error) {
	var zeroLogging zerolog.Logger
	var initErr error

	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			initErr = err
			return
		}

		zeroLogging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger().
			Level(level)
	})

	if initErr != nil {
		return nil, fmt.Errorf("failed to parse log level %q: %w", cfg.Level, initErr)
	}

	return &logger{log: zeroLogging}, nil
}

func (l *logger) Trace(ctx context.Context, obj any) {
	l.log.Trace().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Debug(ctx context.Context, obj any) {
	l.log.Debug().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Info(ctx context.Context, obj any) {
	l.log.Info().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Warn(ctx context.Context, obj any) {
	l.log.Warn().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Error(ctx context.Context, obj any) {
	l.log.Error().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Fatal(ctx context.Context, obj any) {
	l.log.Fatal().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func getCaller(obj any) any {
	switch tr := obj.(type) {
	case error:
		file, line, msg, err := errors.GetCaller(tr)
		obj = operator.Ternary(err != nil, fmt.Sprintf("error cannot get caller, %v", err), fmt.Sprintf("%s:%#v --- %s", file, line, msg))
	case string:
		obj = tr
	default:
		obj = fmt.Sprintf("%#v", tr)
	}

	return obj
}

func getContextFields(ctx context.Context) map[string]any {
	reqstart := appcontext.GetRequestStartTime(ctx)
	timeElapsed := "0ms"
	if !time.Time.IsZero(reqstart) {
		timeElapsed = fmt.Sprintf("%dms", int64(time.Since(reqstart)/time.Millisecond))
	}

	return map[string]any{
		"request_id":      appcontext.GetRequestId(ctx),
		"user_agent":      appcontext.GetUserAgent(ctx),
		"service_version": appcontext.GetServiceVersion(ctx),
		"time_elapsed":    timeElapsed,
	}
}
