package hash

import (
	"strings"
	"testing"
)

func TestNewArgon2(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: []byte("mysecretpassword"),
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: []byte(""),
			wantErr:  false,
		},
		{
			name:     "long password",
			password: []byte(strings.Repeat("a", 1000)),
			wantErr:  false,
		},
		{
			name:     "unicode password",
			password: []byte("пароль密码パスワード"),
			wantErr:  false,
		},
		{
			name:     "password with special chars",
			password: []byte("p@$$w0rd!#%^&*()"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewArgon2(tt.password)
			if tt.wantErr {
				if err == nil {
					t.Error("NewArgon2() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewArgon2() unexpected error: %v", err)
			}
			if !strings.HasPrefix(hash, "$argon2id$v=") {
				t.Errorf("NewArgon2() hash doesn't have argon2id prefix: %s", hash)
			}
			// Hash should have 6 parts separated by $
			parts := strings.Split(hash, "$")
			if len(parts) != 6 {
				t.Errorf("NewArgon2() hash has %d parts, want 6: %s", len(parts), hash)
			}
		})
	}
}

func TestNewArgon2_ProducesDifferentHashes(t *testing.T) {
	password := []byte("samepassword")

	hash1, err := NewArgon2(password)
	if err != nil {
		t.Fatalf("first hash error: %v", err)
	}

	hash2, err := NewArgon2(password)
	if err != nil {
		t.Fatalf("second hash error: %v", err)
	}

	// Due to random salt, same password should produce different hashes
	if hash1 == hash2 {
		t.Error("NewArgon2() produced identical hashes for same password (salt reuse)")
	}
}

func TestCompareArgon2(t *testing.T) {
	// Pre-hash a known password
	password := "testpassword123"
	hash, err := NewArgon2([]byte(password))
	if err != nil {
		t.Fatalf("setup: NewArgon2() error: %v", err)
	}

	tests := []struct {
		name        string
		password    string
		encodedHash string
		wantMatch   bool
		wantErr     bool
		errContains string
	}{
		{
			name:        "correct password matches",
			password:    password,
			encodedHash: hash,
			wantMatch:   true,
			wantErr:     false,
		},
		{
			name:        "wrong password does not match",
			password:    "wrongpassword",
			encodedHash: hash,
			wantMatch:   false,
			wantErr:     false,
		},
		{
			name:        "empty password does not match",
			password:    "",
			encodedHash: hash,
			wantMatch:   false,
			wantErr:     false,
		},
		{
			name:        "invalid hash format - too few parts",
			password:    password,
			encodedHash: "$argon2id$v=19$m=4096",
			wantMatch:   false,
			wantErr:     true,
			errContains: "invalid length of encoded hash",
		},
		{
			name:        "invalid hash format - empty string",
			password:    password,
			encodedHash: "",
			wantMatch:   false,
			wantErr:     true,
			errContains: "invalid length of encoded hash",
		},
		{
			name:        "invalid hash format - wrong version",
			password:    password,
			encodedHash: "$argon2id$v=16$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g",
			wantMatch:   false,
			wantErr:     true,
			errContains: "incompatible argon2 version",
		},
		{
			name:        "invalid hash format - bad base64 salt",
			password:    password,
			encodedHash: "$argon2id$v=19$m=4096,t=3,p=1$!!!invalid!!!$c29tZWhhc2g",
			wantMatch:   false,
			wantErr:     true,
			errContains: "error while decode salt",
		},
		{
			name:        "invalid hash format - bad base64 hash",
			password:    password,
			encodedHash: "$argon2id$v=19$m=4096,t=3,p=1$c29tZXNhbHQ$!!!invalid!!!",
			wantMatch:   false,
			wantErr:     true,
			errContains: "error while decode hash",
		},
		{
			name:        "invalid hash format - bad params",
			password:    password,
			encodedHash: "$argon2id$v=19$invalid_params$c29tZXNhbHQ$c29tZWhhc2g",
			wantMatch:   false,
			wantErr:     true,
			errContains: "error while get values from memory",
		},
		{
			name:        "invalid hash format - bad version format",
			password:    password,
			encodedHash: "$argon2id$v=abc$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g",
			wantMatch:   false,
			wantErr:     true,
			errContains: "error while get argon2 version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CompareArgon2(tt.password, tt.encodedHash)
			if tt.wantErr {
				if err == nil {
					t.Error("CompareArgon2() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CompareArgon2() error = %v, want contains %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompareArgon2() unexpected error: %v", err)
			}
			if match != tt.wantMatch {
				t.Errorf("CompareArgon2() = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestArgon2RoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"with spaces and stuff",
		"p@$$w0rd!#%^&*()",
		strings.Repeat("x", 100),
		"unicode: 你好世界",
	}

	for _, pw := range passwords {
		t.Run("roundtrip_"+pw[:min(len(pw), 15)], func(t *testing.T) {
			hash, err := NewArgon2([]byte(pw))
			if err != nil {
				t.Fatalf("NewArgon2() error: %v", err)
			}

			match, err := CompareArgon2(pw, hash)
			if err != nil {
				t.Fatalf("CompareArgon2() error: %v", err)
			}
			if !match {
				t.Error("CompareArgon2() password should match its own hash")
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
