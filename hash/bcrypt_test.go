package hash

import (
	"strings"
	"testing"
)

func TestCreateBcrypt(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		cost    int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid password with default cost",
			text:    "password123",
			cost:    4,
			wantErr: false,
		},
		{
			name:    "valid password with min cost",
			text:    "password123",
			cost:    4,
			wantErr: false,
		},
		{
			name:    "valid password with higher cost",
			text:    "password123",
			cost:    6,
			wantErr: false,
		},
		{
			name:    "empty password",
			text:    "",
			cost:    10,
			wantErr: false,
		},
		{
			name:    "cost too high",
			text:    "password",
			cost:    35,
			wantErr: true,
			errMsg:  "error while encode bcrypt hash",
		},
		{
			name:    "cost below min is accepted (auto-upgraded)",
			text:    "password",
			cost:    1,
			wantErr: false,
		},
		{
			name:    "password with special characters",
			text:    "p@$$w0rd!#%^&*()",
			cost:    4,
			wantErr: false,
		},
		{
			name:    "unicode password",
			text:    "密码パスワード",
			cost:    4,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := CreateBcrypt(tt.text, tt.cost)
			if tt.wantErr {
				if err == nil {
					t.Error("CreateBcrypt() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("CreateBcrypt() error = %v, want contains %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateBcrypt() unexpected error: %v", err)
			}
			if hash == "" {
				t.Error("CreateBcrypt() returned empty hash")
			}
			// bcrypt hashes start with $2a$ or $2b$
			if !strings.HasPrefix(hash, "$2") {
				t.Errorf("CreateBcrypt() hash doesn't have bcrypt prefix: %s", hash)
			}
		})
	}
}

func TestCreateBcrypt_ProducesDifferentHashes(t *testing.T) {
	password := "samepassword"

	hash1, err := CreateBcrypt(password, 4)
	if err != nil {
		t.Fatalf("first hash error: %v", err)
	}

	hash2, err := CreateBcrypt(password, 4)
	if err != nil {
		t.Fatalf("second hash error: %v", err)
	}

	if hash1 == hash2 {
		t.Error("CreateBcrypt() produced identical hashes for same password (salt reuse)")
	}
}

func TestCompareBcrypt(t *testing.T) {
	password := "testpassword123"
	hash, err := CreateBcrypt(password, 4)
	if err != nil {
		t.Fatalf("setup: CreateBcrypt() error: %v", err)
	}

	tests := []struct {
		name    string
		hash    string
		text    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "correct password matches",
			hash:    hash,
			text:    password,
			wantErr: false,
		},
		{
			name:    "wrong password does not match",
			hash:    hash,
			text:    "wrongpassword",
			wantErr: true,
			errMsg:  "error while compare bcrypt hash",
		},
		{
			name:    "empty password does not match",
			hash:    hash,
			text:    "",
			wantErr: true,
			errMsg:  "error while compare bcrypt hash",
		},
		{
			name:    "invalid hash format - too short",
			hash:    "$2b$10$short",
			text:    password,
			wantErr: true,
			errMsg:  "error while compare bcrypt hash",
		},
		{
			name:    "invalid hash format - empty",
			hash:    "",
			text:    password,
			wantErr: true,
			errMsg:  "error while compare bcrypt hash",
		},
		{
			name:    "invalid hash format - garbage",
			hash:    "not-a-bcrypt-hash-at-all",
			text:    password,
			wantErr: true,
			errMsg:  "error while compare bcrypt hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompareBcrypt(tt.hash, tt.text)
			if tt.wantErr {
				if err == nil {
					t.Error("CompareBcrypt() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("CompareBcrypt() error = %v, want contains %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompareBcrypt() unexpected error: %v", err)
			}
		})
	}
}

func TestBcryptRoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"with spaces",
		"p@$$w0rd!",
		"unicode: 你好",
	}

	for _, pw := range passwords {
		t.Run("roundtrip_"+pw, func(t *testing.T) {
			hash, err := CreateBcrypt(pw, 4)
			if err != nil {
				t.Fatalf("CreateBcrypt() error: %v", err)
			}

			err = CompareBcrypt(hash, pw)
			if err != nil {
				t.Errorf("CompareBcrypt() password should match its own hash: %v", err)
			}
		})
	}
}
