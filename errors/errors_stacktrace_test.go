package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/alpardfm/go-toolkit/codes"
)

func Test_stacktrace_Error(t *testing.T) {
	type fields struct {
		message  string
		cause    error
		code     codes.Code
		file     string
		function string
		line     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "ok",
			fields: fields{message: "failed to start"},
			want:   "Error: failed to start",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stacktrace{
				message:  tt.fields.message,
				cause:    tt.fields.cause,
				code:     tt.fields.code,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			if got := st.Error(); got != tt.want {
				t.Errorf("stacktrace.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stacktrace_ExitCode(t *testing.T) {
	type fields struct {
		message  string
		cause    error
		code     codes.Code
		file     string
		function string
		line     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "no code",
			fields: fields{code: codes.NoCode},
			want:   1,
		},
		{
			name:   "auth failuer",
			fields: fields{code: codes.CodeAuth},
			want:   int(codes.CodeAuth),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stacktrace{
				message:  tt.fields.message,
				cause:    tt.fields.cause,
				code:     tt.fields.code,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			if got := st.ExitCode(); got != tt.want {
				t.Errorf("stacktrace.ExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stacktrace_Unwrap(t *testing.T) {
	t.Run("returns cause when cause is set", func(t *testing.T) {
		cause := fmt.Errorf("root cause")
		st := &stacktrace{
			message: "wrapped error",
			cause:   cause,
			code:    codes.NoCode,
		}
		if got := st.Unwrap(); got != cause {
			t.Errorf("stacktrace.Unwrap() = %v, want %v", got, cause)
		}
	})

	t.Run("returns nil when cause is nil", func(t *testing.T) {
		st := &stacktrace{
			message: "no cause",
			cause:   nil,
			code:    codes.NoCode,
		}
		if got := st.Unwrap(); got != nil {
			t.Errorf("stacktrace.Unwrap() = %v, want nil", got)
		}
	})

	t.Run("errors.Is traverses single-level chain", func(t *testing.T) {
		target := fmt.Errorf("sentinel error")
		st := &stacktrace{
			message: "wrapped",
			cause:   target,
			code:    codes.NoCode,
		}
		if !errors.Is(st, target) {
			t.Errorf("errors.Is(st, target) = false, want true")
		}
	})

	t.Run("errors.Is traverses multi-level chain", func(t *testing.T) {
		target := fmt.Errorf("deep sentinel")
		inner := &stacktrace{
			message: "inner wrap",
			cause:   target,
			code:    codes.NoCode,
		}
		outer := &stacktrace{
			message: "outer wrap",
			cause:   inner,
			code:    codes.NoCode,
		}
		if !errors.Is(outer, target) {
			t.Errorf("errors.Is(outer, target) = false, want true for multi-level chain")
		}
	})

	t.Run("errors.As extracts typed error from chain", func(t *testing.T) {
		target := &typedTestError{detail: "test detail"}
		st := &stacktrace{
			message: "wrapped typed",
			cause:   target,
			code:    codes.NoCode,
		}
		var extracted *typedTestError
		if !errors.As(st, &extracted) {
			t.Errorf("errors.As(st, &extracted) = false, want true")
		}
		if extracted.detail != "test detail" {
			t.Errorf("extracted.detail = %q, want %q", extracted.detail, "test detail")
		}
	})
}

// typedTestError is a helper type for testing errors.As
type typedTestError struct {
	detail string
}

func (e *typedTestError) Error() string {
	return e.detail
}
