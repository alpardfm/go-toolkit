package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// NewSHA256WithKey generates a HMAC SHA256 hash with a given key.
func NewSHA256WithKey(text, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	io.WriteString(hash, text)
	return hex.EncodeToString(hash.Sum(nil))
}

// NewSHA256 generates a SHA256 hash.
func NewSHA256(text string) string {
	hash := sha256.New()
	io.WriteString(hash, text)
	return hex.EncodeToString(hash.Sum(nil))
}
