package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/language"
)

// currentTestFile returns the absolute path of the current test file using runtime.Caller.
func currentTestFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

func TestApp_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "OK",
			want: "Error: invalid format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &App{
				sys: NewWithCode(codes.CodeBadRequest, "invalid format"),
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("App.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompile(t *testing.T) {
	type args struct {
		err  error
		lang string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 App
	}{
		{
			name: "ok",
			args: args{err: NewWithCode(codes.CodeAuthFailure, "auth failed"), lang: language.English},
			want: http.StatusUnauthorized,
			want1: App{
				Code:  codes.CodeAuthFailure,
				Title: "Unauthorized",
				Body:  "Unauthorized access. You are not authorized to access this resource.",
				sys:   NewWithCode(codes.CodeAuthFailure, "auth failed"),
			},
		},
		{
			name: "not ok",
			args: args{err: NewWithCode(codes.NoCode, "no code"), lang: language.English},
			want: http.StatusInternalServerError,
			want1: App{
				Code:  codes.NoCode,
				Title: "Service Error Not Defined",
				Body:  "Unknown error. Please contact admin",
				sys:   NewWithCode(codes.CodeAuthFailure, "auth failed"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Compile(tt.args.err, tt.args.lang)
			if got != tt.want {
				t.Errorf("Compile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1.Code, tt.want1.Code) {
				t.Errorf("Compile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetCaller(t *testing.T) {
	testFile := currentTestFile()

	t.Run("ok", func(t *testing.T) {
		// Create the error and record the line number where it's created.
		// runtime.Caller(2) inside create() will point to this call site.
		err := NewWithCode(codes.CodeBadRequest, "bad request")
		_, _, expectedLine, _ := runtime.Caller(0)
		expectedLine-- // the previous line is where the error was created

		got, got1, got2, gotErr := GetCaller(err)
		if gotErr != nil {
			t.Errorf("GetCaller() error = %v, wantErr false", gotErr)
			return
		}
		if !strings.HasSuffix(got, "errors/errors_test.go") {
			t.Errorf("GetCaller() got file = %v, want suffix %v", got, "errors/errors_test.go")
		}
		if got != testFile {
			t.Errorf("GetCaller() got file = %v, want %v", got, testFile)
		}
		if got1 != expectedLine {
			t.Errorf("GetCaller() got line = %v, want %v", got1, expectedLine)
		}
		if got2 != "bad request" {
			t.Errorf("GetCaller() got message = %v, want %v", got2, "bad request")
		}
	})

	t.Run("not ok - non-stacktrace error", func(t *testing.T) {
		got, got1, got2, gotErr := GetCaller(fmt.Errorf(""))
		if gotErr == nil {
			t.Errorf("GetCaller() error = nil, wantErr true")
			return
		}
		if got != "" {
			t.Errorf("GetCaller() got = %v, want empty string", got)
		}
		if got1 != 0 {
			t.Errorf("GetCaller() got1 = %v, want 0", got1)
		}
		if got2 != "" {
			t.Errorf("GetCaller() got2 = %v, want empty string", got2)
		}
	})
}
