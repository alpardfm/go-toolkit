package appcontext

import (
	"context"
	"reflect"
	"testing"

	"github.com/alpardfm/go-toolkit/language"
)

func TestSetAcceptLanguage(t *testing.T) {
	type args struct {
		ctx      context.Context
		language string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), language: language.English},
			want: context.WithValue(context.Background(), acceptLanguage, language.English),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAcceptLanguage(tt.args.ctx, tt.args.language); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAcceptLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAcceptLanguage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: language.English,
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), acceptLanguage, language.Indonesian)},
			want: language.Indonesian,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAcceptLanguage(tt.args.ctx); got != tt.want {
				t.Errorf("GetAcceptLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestId(t *testing.T) {
	type args struct {
		ctx context.Context
		rid string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), rid: "randomized-request-id"},
			want: context.WithValue(context.Background(), requestId, "randomized-request-id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestId(tt.args.ctx, tt.args.rid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestId(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestId, "randomized-request-id")},
			want: "randomized-request-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestId(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetServiceVersion(t *testing.T) {
	type args struct {
		ctx     context.Context
		version string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), version: "v1.0.0"},
			want: context.WithValue(context.Background(), serviceVersion, "v1.0.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetServiceVersion(tt.args.ctx, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetServiceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceVersion(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), serviceVersion, "v1.0.0")},
			want: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceVersion(tt.args.ctx); got != tt.want {
				t.Errorf("GetServiceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUserAgent(t *testing.T) {
	type args struct {
		ctx context.Context
		ua  string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), ua: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)"},
			want: context.WithValue(context.Background(), userAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserAgent(tt.args.ctx, tt.args.ua); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserAgent(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), userAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")},
			want: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserAgent(tt.args.ctx); got != tt.want {
				t.Errorf("GetUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
