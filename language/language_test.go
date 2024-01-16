package language

import (
	"net/http"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	type param struct {
		lang string
		code int
	}
	tests := []struct {
		name  string
		param param
		want  string
	}{
		{
			name:  "Switch to english language",
			param: param{lang: English, code: http.StatusContinue},
			want:  "Continue",
		},
		{
			name:  "Switch to indonesian language",
			param: param{lang: Indonesian, code: http.StatusContinue},
			want:  "Lanjutkan",
		},
		{
			name:  "Unknown language",
			param: param{lang: "Alien", code: http.StatusContinue},
			want:  "Continue",
		},
		{
			name:  "Unknown code",
			param: param{lang: English, code: 1},
			want:  "Response Unknown",
		},
		{
			name:  "Unknown language and code",
			param: param{lang: "Alien", code: 1},
			want:  "Response Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTTPStatusText(tt.param.lang, tt.param.code); got != tt.want {
				t.Errorf("HTTStatusText() = %v, want %v", got, tt.want)
			}
		})
	}
}
