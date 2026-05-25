package hash

import (
	"testing"
)

func BenchmarkNewArgon2(b *testing.B) {
	b.ReportAllocs()
	password := []byte("correct-horse-battery-staple")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = NewArgon2(password)
	}
}

func BenchmarkCompareArgon2(b *testing.B) {
	b.ReportAllocs()
	password := "correct-horse-battery-staple"
	encodedHash, err := NewArgon2([]byte(password))
	if err != nil {
		b.Fatalf("failed to create argon2 hash: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = CompareArgon2(password, encodedHash)
	}
}

func BenchmarkCreateBcrypt(b *testing.B) {
	b.ReportAllocs()
	plainText := "correct-horse-battery-staple"
	cost := 10
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = CreateBcrypt(plainText, cost)
	}
}

func BenchmarkCompareBcrypt(b *testing.B) {
	b.ReportAllocs()
	plainText := "correct-horse-battery-staple"
	hashed, err := CreateBcrypt(plainText, 10)
	if err != nil {
		b.Fatalf("failed to create bcrypt hash: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CompareBcrypt(hashed, plainText)
	}
}

func BenchmarkNewSHA256(b *testing.B) {
	b.ReportAllocs()
	text := "the quick brown fox jumps over the lazy dog"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewSHA256(text)
	}
}

func BenchmarkNewSHA256WithKey(b *testing.B) {
	b.ReportAllocs()
	text := "the quick brown fox jumps over the lazy dog"
	key := "my-secret-hmac-key-256-bits-long"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewSHA256WithKey(text, key)
	}
}
