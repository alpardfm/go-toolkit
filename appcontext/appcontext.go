package appcontext

import (
	"context"
	"time"

	"github.com/alpardfm/go-toolkit/language"
)

type contextKey string

const (
	acceptLanguage   contextKey = "AcceptLanguage"
	requestId        contextKey = "RequestId"
	serviceVersion   contextKey = "ServiceVersion"
	userAgent        contextKey = "UserAgent"
	requestStartTime contextKey = "RequestStartTime"
)

func SetAcceptLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, acceptLanguage, lang)
}

func GetAcceptLanguage(ctx context.Context) string {
	lang, ok := ctx.Value(acceptLanguage).(string)
	if !ok {
		// return english as the default language
		return language.English
	}

	return lang
}

func SetRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestId, rid)
}

func GetRequestId(ctx context.Context) string {
	rid, ok := ctx.Value(requestId).(string)
	if !ok {
		return ""
	}
	return rid
}

func SetServiceVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, serviceVersion, version)
}

func GetServiceVersion(ctx context.Context) string {
	version, ok := ctx.Value(serviceVersion).(string)
	if !ok {
		return ""
	}
	return version
}

func SetUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, userAgent, ua)
}

func GetUserAgent(ctx context.Context) string {
	ua, ok := ctx.Value(userAgent).(string)
	if !ok {
		return ""
	}
	return ua
}

func SetRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTime, t)
}

func GetRequestStartTime(ctx context.Context) time.Time {
	t, _ := ctx.Value(requestStartTime).(time.Time)
	return t
}
