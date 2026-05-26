package encrypt

import (
	"testing"
)

// FuzzDecryptAES256GCM tests that DecryptAES256GCM never panics on arbitrary input.
func FuzzDecryptAES256GCM(f *testing.F) {
	// Seed corpus
	f.Add([]byte("not-base64"), []byte("12345678901234567890123456789012"))
	f.Add([]byte(""), []byte("12345678901234567890123456789012"))
	f.Add([]byte("aGVsbG8="), []byte("short"))
	f.Add([]byte("dGVzdA=="), []byte("12345678901234567890123456789012"))

	// Add a valid encrypted value
	key := make([]byte, 32)
	copy(key, []byte("this-is-a-32-byte-key-for-aes!!"))
	encrypted, _ := EncryptAES256GCM([]byte("test"), key)
	if encrypted != nil {
		f.Add(encrypted, key)
	}

	f.Fuzz(func(t *testing.T, ciphertext, key []byte) {
		// Should never panic, errors are acceptable
		_, _ = DecryptAES256GCM(ciphertext, key)
	})
}

// FuzzEncryptAES256GCM tests that EncryptAES256GCM never panics on arbitrary input.
func FuzzEncryptAES256GCM(f *testing.F) {
	f.Add([]byte("hello"), []byte("12345678901234567890123456789012"))
	f.Add([]byte(""), []byte("12345678901234567890123456789012"))
	f.Add([]byte("test"), []byte("short"))
	f.Add([]byte("data"), []byte(""))

	f.Fuzz(func(t *testing.T, plaintext, key []byte) {
		// Should never panic, errors are acceptable
		result, err := EncryptAES256GCM(plaintext, key)
		if err != nil {
			return
		}

		// If encryption succeeded, decryption with same key should work
		decrypted, err := DecryptAES256GCM(result, key)
		if err != nil {
			t.Errorf("encrypt succeeded but decrypt failed: %v", err)
			return
		}
		if string(decrypted) != string(plaintext) {
			t.Errorf("round trip failed: got %q, want %q", decrypted, plaintext)
		}
	})
}
