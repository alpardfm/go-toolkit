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
	"github.com/rs/zerolog/log"
)

var once = sync.Once{}

type Interface interface {
	// TODO add Debugf
	Trace(ctx context.Context, obj interface{})
	Debug(ctx context.Context, obj interface{})
	Info(ctx context.Context, obj interface{})
	Warn(ctx context.Context, obj interface{})
	Error(ctx context.Context, obj interface{})
	Fatal(ctx context.Context, obj interface{})
}

type Config struct {
	Level string
}

type logger struct {
	log zerolog.Logger
}

func Init(cfg Config) Interface {
	var zeroLogging zerolog.Logger
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("failed to parse error level with err: %v", err))
		}

		zeroLogging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3). //Hard code to 3 for now.
			Logger().
			Level(level)
	})

	return &logger{log: zeroLogging}
}

func (l *logger) Trace(ctx context.Context, obj interface{}) {
	l.log.Trace().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Debug(ctx context.Context, obj interface{}) {
	l.log.Debug().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Info(ctx context.Context, obj interface{}) {
	l.log.Info().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Warn(ctx context.Context, obj interface{}) {
	l.log.Warn().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Error(ctx context.Context, obj interface{}) {
	l.log.Error().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Fatal(ctx context.Context, obj interface{}) {
	l.log.Fatal().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func getCaller(obj interface{}) interface{} {
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

func getContextFields(ctx context.Context) map[string]interface{} {
	reqstart := appcontext.GetRequestStartTime(ctx)
	timeElapsed := "0ms"
	if !time.Time.IsZero(reqstart) {
		timeElapsed = fmt.Sprintf("%dms", int64(time.Since(reqstart)/time.Millisecond))
	}

	return map[string]interface{}{
		"request_id":      appcontext.GetRequestId(ctx),
		"user_agent":      appcontext.GetUserAgent(ctx),
		"service_version": appcontext.GetServiceVersion(ctx),
		"time_elapsed":    timeElapsed,
	}
}
