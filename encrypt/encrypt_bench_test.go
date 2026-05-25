package encrypt

import (
	"crypto/rand"
	"testing"
)

func BenchmarkEncryptAES256GCM(b *testing.B) {
	b.ReportAllocs()
	plaintext := []byte("sensitive data that needs encryption for benchmarking purposes")
	key := make([]byte, 32) // AES-256 requires 32-byte key
	if _, err := rand.Read(key); err != nil {
		b.Fatalf("failed to generate key: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = EncryptAES256GCM(plaintext, key)
	}
}

func BenchmarkDecryptAES256GCM(b *testing.B) {
	b.ReportAllocs()
	plaintext := []byte("sensitive data that needs encryption for benchmarking purposes")
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		b.Fatalf("failed to generate key: %v", err)
	}
	ciphertext, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		b.Fatalf("failed to encrypt: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = DecryptAES256GC(ciphertext, key)
	}
}
