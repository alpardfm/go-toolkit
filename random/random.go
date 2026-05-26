package random

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/google/uuid"
)

// GenerateInt generates a cryptographically secure random positive integer with the specified number of digits.
// Each digit is guaranteed to be non-zero (1-9).
// Uses crypto/rand for secure randomness suitable for tokens, OTPs, and security-sensitive use cases.
func GenerateInt(length int) (int64, error) {
	if length <= 0 {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, "length must be greater than 0")
	}
	if length > 18 {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, "length must be 18 or less to fit in int64")
	}

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		// Generate a random number in range [1, 9]
		n, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return 0, errors.NewWithCode(codes.CodeInvalidValue, "failed to generate random number: %v", err)
		}
		// n is [0, 8], add 1 to get [1, 9]
		digit := n.Int64() + 1
		sb.WriteByte(byte('0' + digit))
	}

	// Parse the string to int64
	var result int64
	for _, c := range sb.String() {
		result = result*10 + int64(c-'0')
	}

	return result, nil
}

// GeneratePUID generates a pseudo-unique identifier of the specified length
// by combining the given MSISDN with a UUID and hashing via SHA1.
func GeneratePUID(msisdn string, length int) (string, error) {
	if length <= 0 {
		return "", errors.NewWithCode(codes.CodeInvalidValue, "length must be greater than 0")
	}
	if length > 32 {
		return "", errors.NewWithCode(codes.CodeInvalidValue, "length must be 32 or less (UUID hex length without dashes)")
	}

	combinedString := msisdn + uuid.New().String()
	puid := uuid.NewSHA1(uuid.Nil, []byte(combinedString))

	result := strings.ReplaceAll(puid.String(), "-", "")[0:length]
	return result, nil
}
