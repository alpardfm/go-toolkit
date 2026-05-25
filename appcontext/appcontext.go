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

// SetAcceptLanguage stores the accept language value in the context.
func SetAcceptLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, acceptLanguage, lang)
}

// GetAcceptLanguage retrieves the accept language from the context.
// Returns English as the default if not set.
func GetAcceptLanguage(ctx context.Context) string {
	lang, ok := ctx.Value(acceptLanguage).(string)
	if !ok {
		// return english as the default language
		return language.English
	}

	return lang
}

// SetRequestId stores the request ID in the context.
func SetRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestId, rid)
}

// GetRequestId retrieves the request ID from the context.
// Returns an empty string if not set.
func GetRequestId(ctx context.Context) string {
	rid, ok := ctx.Value(requestId).(string)
	if !ok {
		return ""
	}
	return rid
}

// SetServiceVersion stores the service version in the context.
func SetServiceVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, serviceVersion, version)
}

// GetServiceVersion retrieves the service version from the context.
// Returns an empty string if not set.
func GetServiceVersion(ctx context.Context) string {
	version, ok := ctx.Value(serviceVersion).(string)
	if !ok {
		return ""
	}
	return version
}

// SetUserAgent stores the user agent string in the context.
func SetUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, userAgent, ua)
}

// GetUserAgent retrieves the user agent string from the context.
// Returns an empty string if not set.
func GetUserAgent(ctx context.Context) string {
	ua, ok := ctx.Value(userAgent).(string)
	if !ok {
		return ""
	}
	return ua
}

// SetRequestStartTime stores the request start time in the context.
func SetRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTime, t)
}

// GetRequestStartTime retrieves the request start time from the context.
// Returns a zero time if not set.
func GetRequestStartTime(ctx context.Context) time.Time {
	t, _ := ctx.Value(requestStartTime).(time.Time)
	return t
}
