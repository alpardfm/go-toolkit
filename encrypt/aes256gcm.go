package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

func EncryptAES256GCM(text, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to create new block cipher, %v", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to create new ciper gcm, %v", err))
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to read random value for nonce, %v", err))
	}

	s := string(gcm.Seal(nonce, nonce, []byte(text), nil))
	s = hex.EncodeToString([]byte(s))
	s = base64.StdEncoding.EncodeToString([]byte(s))

	return []byte(s), nil
}

func DecryptAES256GC(text, key []byte) ([]byte, error) {
	text, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to decode base64 string, %v", err))
	}

	text, err = hex.DecodeString(string(text))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to decode hex string, %v", err))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to create new block cipher, %v", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("failed to create new cipher gcm, %v", err))
	}

	nonceSize := gcm.NonceSize()
	nonce, encrypedText := text[:nonceSize], text[nonceSize:]

	if res, err := gcm.Open(nil, nonce, encrypedText, nil); err != nil {
		return nil, errors.NewWithCode(codes.CodeAES256GCMOpenError, fmt.Sprintf("failed to open encrypted aes-256-gcm, %v", err))
	} else {
		return res, nil
	}
}
