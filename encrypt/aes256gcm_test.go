package encrypt

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestEncryptAES256GCM(t *testing.T) {
	// Valid 32-byte key (AES-256 requires exactly 32 bytes)
	key := make([]byte, 32)
	copy(key, []byte("this-is-a-32-byte-key-for-aes!!"))

	tests := []struct {
		name    string
		text    []byte
		key     []byte
		wantErr bool
		errMsg  string
	}{
		{
			name:    "encrypt valid plaintext",
			text:    []byte("hello world"),
			key:     key,
			wantErr: false,
		},
		{
			name:    "encrypt empty plaintext",
			text:    []byte(""),
			key:     key,
			wantErr: false,
		},
		{
			name:    "encrypt long plaintext",
			text:    []byte(strings.Repeat("a", 10000)),
			key:     key,
			wantErr: false,
		},
		{
			name:    "encrypt with invalid key length (too short)",
			text:    []byte("hello"),
			key:     []byte("short"),
			wantErr: true,
			errMsg:  "failed to create new block cipher",
		},
		{
			name:    "encrypt with invalid key length (too long)",
			text:    []byte("hello"),
			key:     []byte(strings.Repeat("k", 64)),
			wantErr: true,
			errMsg:  "failed to create new block cipher",
		},
		{
			name:    "encrypt with 16-byte key (AES-128)",
			text:    []byte("hello"),
			key:     []byte("sixteen-byte-key"),
			wantErr: false, // AES supports 16, 24, 32 byte keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncryptAES256GCM(tt.text, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Errorf("EncryptAES256GCM() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("EncryptAES256GCM() error = %v, want contains %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("EncryptAES256GCM() unexpected error: %v", err)
			}
			if len(result) == 0 {
				t.Error("EncryptAES256GCM() returned empty result")
			}
			// Result should be valid base64
			if _, err := base64.StdEncoding.DecodeString(string(result)); err != nil {
				t.Errorf("EncryptAES256GCM() result is not valid base64: %v", err)
			}
			// Encrypted text should differ from plaintext
			if string(result) == string(tt.text) {
				t.Error("EncryptAES256GCM() encrypted text equals plaintext")
			}
		})
	}
}

func TestDecryptAES256GCM(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("this-is-a-32-byte-key-for-aes!!"))

	tests := []struct {
		name    string
		text    []byte
		key     []byte
		wantErr bool
		errMsg  string
	}{
		{
			name:    "decrypt with invalid base64",
			text:    []byte("not-valid-base64!!!"),
			key:     key,
			wantErr: true,
			errMsg:  "failed to decode base64 string",
		},
		{
			name:    "decrypt with valid base64 but invalid hex",
			text:    []byte(base64.StdEncoding.EncodeToString([]byte("not-hex"))),
			key:     key,
			wantErr: true,
			errMsg:  "failed to decode hex string",
		},
		{
			name:    "decrypt with wrong key",
			text:    nil, // will be set dynamically
			key:     []byte("wrong-key-that-is-32-bytes-long!"),
			wantErr: true,
			errMsg:  "failed to open encrypted aes-256-gcm",
		},
		{
			name:    "decrypt with short key",
			text:    nil, // will be set dynamically
			key:     []byte("short"),
			wantErr: true,
			errMsg:  "failed to create new block cipher",
		},
	}

	// Pre-encrypt a value for tests that need valid ciphertext
	validCiphertext, err := EncryptAES256GCM([]byte("test-data"), key)
	if err != nil {
		t.Fatalf("setup: failed to encrypt test data: %v", err)
	}

	for i := range tests {
		if tests[i].text == nil {
			tests[i].text = validCiphertext
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptAES256GCM(tt.text, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Errorf("DecryptAES256GCM() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("DecryptAES256GCM() error = %v, want contains %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("DecryptAES256GCM() unexpected error: %v", err)
			}
		})
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("this-is-a-32-byte-key-for-aes!!"))

	plaintexts := []string{
		"hello world",
		"",
		"short",
		strings.Repeat("long text ", 1000),
		"special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		"unicode: 你好世界 🌍🚀",
	}

	for _, pt := range plaintexts {
		t.Run("roundtrip_"+pt[:min(len(pt), 20)], func(t *testing.T) {
			encrypted, err := EncryptAES256GCM([]byte(pt), key)
			if err != nil {
				t.Fatalf("EncryptAES256GCM() error: %v", err)
			}

			decrypted, err := DecryptAES256GCM(encrypted, key)
			if err != nil {
				t.Fatalf("DecryptAES256GCM() error: %v", err)
			}

			if string(decrypted) != pt {
				t.Errorf("round trip failed: got %q, want %q", string(decrypted), pt)
			}
		})
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("this-is-a-32-byte-key-for-aes!!"))

	plaintext := []byte("same input")

	ct1, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		t.Fatalf("first encrypt error: %v", err)
	}

	ct2, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		t.Fatalf("second encrypt error: %v", err)
	}

	// Due to random nonce, same plaintext should produce different ciphertexts
	if string(ct1) == string(ct2) {
		t.Error("EncryptAES256GCM() produced identical ciphertexts for same plaintext (nonce reuse)")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
