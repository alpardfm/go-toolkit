package log

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid info level",
			cfg:     Config{Level: "info"},
			wantErr: false,
		},
		{
			name:    "valid debug level",
			cfg:     Config{Level: "debug"},
			wantErr: false,
		},
		{
			name:    "valid error level",
			cfg:     Config{Level: "error"},
			wantErr: false,
		},
		{
			name:    "valid warn level",
			cfg:     Config{Level: "warn"},
			wantErr: false,
		},
		{
			name:    "valid trace level",
			cfg:     Config{Level: "trace"},
			wantErr: false,
		},
		{
			name:    "invalid level",
			cfg:     Config{Level: "invalid"},
			wantErr: true,
			errMsg:  "failed to parse log level",
		},
		{
			name:    "empty level defaults to no-level (valid in zerolog)",
			cfg:     Config{Level: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := Init(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Error("Init() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Init() error = %v, want contains %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("Init() unexpected error: %v", err)
			}
			if l == nil {
				t.Error("Init() returned nil logger")
			}
		})
	}
}

func TestInit_MultipleInstances(t *testing.T) {
	// Verify that Init can be called multiple times with different configs
	l1, err := Init(Config{Level: "info"})
	if err != nil {
		t.Fatalf("first Init() error: %v", err)
	}

	l2, err := Init(Config{Level: "debug"})
	if err != nil {
		t.Fatalf("second Init() error: %v", err)
	}

	if l1 == l2 {
		t.Error("Init() returned same instance for different configs")
	}
}

func TestInit_CustomWriter(t *testing.T) {
	var buf bytes.Buffer

	l, err := Init(Config{Level: "info", Writer: &buf})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	l.Info(context.Background(), "test message")

	output := buf.String()
	if output == "" {
		t.Error("expected log output in buffer, got empty")
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer

	// Set level to info — trace and debug should be suppressed
	l, err := Init(Config{Level: "info", Writer: &buf})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	l.Trace(context.Background(), "trace msg")
	l.Debug(context.Background(), "debug msg")

	output := buf.String()
	if strings.Contains(output, "trace msg") {
		t.Error("trace message should be suppressed at info level")
	}
	if strings.Contains(output, "debug msg") {
		t.Error("debug message should be suppressed at info level")
	}

	// Info and above should appear
	buf.Reset()
	l.Info(context.Background(), "info msg")
	l.Warn(context.Background(), "warn msg")
	l.Error(context.Background(), "error msg")

	output = buf.String()
	if !strings.Contains(output, "info msg") {
		t.Error("info message should appear at info level")
	}
	if !strings.Contains(output, "warn msg") {
		t.Error("warn message should appear at info level")
	}
	if !strings.Contains(output, "error msg") {
		t.Error("error message should appear at info level")
	}
}
