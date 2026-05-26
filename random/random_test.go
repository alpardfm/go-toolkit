package random

import (
	"fmt"
	"testing"
)

func TestGenerateInt(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "generate 1 digit",
			length:  1,
			wantErr: false,
		},
		{
			name:    "generate 6 digits",
			length:  6,
			wantErr: false,
		},
		{
			name:    "generate 10 digits",
			length:  10,
			wantErr: false,
		},
		{
			name:    "generate 18 digits (max)",
			length:  18,
			wantErr: false,
		},
		{
			name:    "zero length returns error",
			length:  0,
			wantErr: true,
			errMsg:  "length must be greater than 0",
		},
		{
			name:    "negative length returns error",
			length:  -1,
			wantErr: true,
			errMsg:  "length must be greater than 0",
		},
		{
			name:    "length exceeds int64 capacity",
			length:  19,
			wantErr: true,
			errMsg:  "length must be 18 or less",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateInt(tt.length)
			if tt.wantErr {
				if err == nil {
					t.Error("GenerateInt() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GenerateInt() unexpected error: %v", err)
			}
			if result <= 0 {
				t.Errorf("GenerateInt() result = %d, want positive", result)
			}

			// Verify digit count
			s := fmt.Sprintf("%d", result)
			if len(s) != tt.length {
				t.Errorf("GenerateInt() result has %d digits, want %d", len(s), tt.length)
			}

			// Verify no zeros
			for _, c := range s {
				if c == '0' {
					t.Errorf("GenerateInt() result contains zero: %s", s)
					break
				}
			}
		})
	}
}

func TestGenerateInt_Uniqueness(t *testing.T) {
	// Generate multiple values and verify they're not all the same
	results := make(map[int64]bool)
	for i := 0; i < 100; i++ {
		val, err := GenerateInt(6)
		if err != nil {
			t.Fatalf("GenerateInt() error: %v", err)
		}
		results[val] = true
	}

	// With 6 digits (1-9 each), we should get many unique values
	if len(results) < 50 {
		t.Errorf("GenerateInt() produced only %d unique values out of 100, expected more randomness", len(results))
	}
}

func TestGeneratePUID(t *testing.T) {
	tests := []struct {
		name    string
		msisdn  string
		length  int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid MSISDN with length 8",
			msisdn:  "+6281234567890",
			length:  8,
			wantErr: false,
		},
		{
			name:    "valid MSISDN with length 16",
			msisdn:  "+6281234567890",
			length:  16,
			wantErr: false,
		},
		{
			name:    "valid MSISDN with max length 32",
			msisdn:  "+6281234567890",
			length:  32,
			wantErr: false,
		},
		{
			name:    "empty MSISDN",
			msisdn:  "",
			length:  8,
			wantErr: false,
		},
		{
			name:    "zero length returns error",
			msisdn:  "+6281234567890",
			length:  0,
			wantErr: true,
			errMsg:  "length must be greater than 0",
		},
		{
			name:    "negative length returns error",
			msisdn:  "+6281234567890",
			length:  -1,
			wantErr: true,
			errMsg:  "length must be greater than 0",
		},
		{
			name:    "length exceeds max returns error",
			msisdn:  "+6281234567890",
			length:  33,
			wantErr: true,
			errMsg:  "length must be 32 or less",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GeneratePUID(tt.msisdn, tt.length)
			if tt.wantErr {
				if err == nil {
					t.Error("GeneratePUID() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GeneratePUID() unexpected error: %v", err)
			}
			if len(result) != tt.length {
				t.Errorf("GeneratePUID() result length = %d, want %d", len(result), tt.length)
			}
		})
	}
}

func TestGeneratePUID_Uniqueness(t *testing.T) {
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		val, err := GeneratePUID("+6281234567890", 16)
		if err != nil {
			t.Fatalf("GeneratePUID() error: %v", err)
		}
		results[val] = true
	}

	// Should produce unique values due to UUID randomness
	if len(results) < 90 {
		t.Errorf("GeneratePUID() produced only %d unique values out of 100", len(results))
	}
}
