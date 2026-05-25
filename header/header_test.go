package header

import (
	"testing"
)

func TestHeaderKeyConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{
			name:     "KeyRequestID",
			constant: KeyRequestID,
			want:     "x-request-id",
		},
		{
			name:     "KeyAuthorization",
			constant: KeyAuthorization,
			want:     "authorization",
		},
		{
			name:     "KeyUserAgent",
			constant: KeyUserAgent,
			want:     "user-agent",
		},
		{
			name:     "KeyContentType",
			constant: KeyContentType,
			want:     "content-type",
		},
		{
			name:     "KeyContentAccept",
			constant: KeyContentAccept,
			want:     "accept",
		},
		{
			name:     "KeyAcceptLanguage",
			constant: KeyAcceptLanguage,
			want:     "accept-language",
		},
		{
			name:     "KeyCacheControl",
			constant: KeyCacheControl,
			want:     "cache-control",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}

func TestContentTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{
			name:     "ContentTypeJSON",
			constant: ContentTypeJSON,
			want:     "application/json",
		},
		{
			name:     "ContentTypeXML",
			constant: ContentTypeXML,
			want:     "application/xml",
		},
		{
			name:     "ContentTypeForm",
			constant: ContentTypeForm,
			want:     "application/x-www-form-urlencoded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}

func TestMediaTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{
			name:     "MediaTextPlain",
			constant: MediaTextPlain,
			want:     "text/plain",
		},
		{
			name:     "MediaTextHTML",
			constant: MediaTextHTML,
			want:     "text/html",
		},
		{
			name:     "MediaTextCSV",
			constant: MediaTextCSV,
			want:     "text/csv",
		},
		{
			name:     "MediaTextXML",
			constant: MediaTextXML,
			want:     "text/xml",
		},
		{
			name:     "MediaImageGIF",
			constant: MediaImageGIF,
			want:     "image/gif",
		},
		{
			name:     "MediaImageJPEG",
			constant: MediaImageJPEG,
			want:     "image/jpeg",
		},
		{
			name:     "MediaImagePNG",
			constant: MediaImagePNG,
			want:     "image/png",
		},
		{
			name:     "MediaImageWEBP",
			constant: MediaImageWEBP,
			want:     "image/webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}

func TestCacheControlConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{
			name:     "CacheControlNoCache",
			constant: CacheControlNoCache,
			want:     "no-cache",
		},
		{
			name:     "CacheControlNoStore",
			constant: CacheControlNoStore,
			want:     "no-store",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}
