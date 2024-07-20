package hash

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateBcrypt(plainText string, cost int) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainText), cost)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeBcryptEncodeHashError, fmt.Sprintf("error while encode bcrypt hash, %s", err.Error()))
	}

	return string(hashed), nil
}

func CompareBcrypt(hashedText string, plainText string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	if err != nil {
		return errors.NewWithCode(codes.CodeBcryptCompareHashError, fmt.Sprintf("error while compare bcrypt hash, %s", err.Error()))
	}

	return nil
}
