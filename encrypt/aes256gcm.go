package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

// EncryptAES256GCM encrypts plaintext using AES-256-GCM with the provided 32-byte key.
// The result is base64-encoded and includes the nonce prepended to the ciphertext.
func EncryptAES256GCM(text, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to create new block cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to create new cipher gcm: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to read random value for nonce: %v", err)
	}

	s := string(gcm.Seal(nonce, nonce, []byte(text), nil))
	s = hex.EncodeToString([]byte(s))
	s = base64.StdEncoding.EncodeToString([]byte(s))

	return []byte(s), nil
}

// DecryptAES256GCM decrypts ciphertext that was encrypted with EncryptAES256GCM using the same 32-byte key.
func DecryptAES256GCM(text, key []byte) ([]byte, error) {
	text, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to decode base64 string: %v", err)
	}

	text, err = hex.DecodeString(string(text))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to decode hex string: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to create new block cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, "failed to create new cipher gcm: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		return nil, errors.NewWithCode(codes.CodeAES256GCMOpenError, "ciphertext too short: must be at least %d bytes", nonceSize)
	}
	nonce, encrypedText := text[:nonceSize], text[nonceSize:]

	if res, err := gcm.Open(nil, nonce, encrypedText, nil); err != nil {
		return nil, errors.NewWithCode(codes.CodeAES256GCMOpenError, "failed to open encrypted aes-256-gcm: %v", err)
	} else {
		return res, nil
	}
}
