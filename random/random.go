package random

import (
	"math/rand"
	"strings"
	"time"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/convert"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/google/uuid"
)

// GenerateInt generates a random positive integer with the specified number of digits.
// Each digit is guaranteed to be non-zero (1-9).
func GenerateInt(length int) (int64, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		s, err := convert.ToString(rng.Intn(9))
		if err != nil {
			return 0, errors.NewWithCode(codes.CodeInvalidValue, err.Error())
		}

		if s == "0" {
			i -= 1
			continue
		}

		sb.WriteString(s)
	}

	iValue, err := convert.ToInt64(sb.String())
	if err != nil {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, err.Error())
	}

	return iValue, nil
}

// GeneratePUID generates a pseudo-unique identifier of the specified length
// by combining the given MSISDN with a UUID and hashing via SHA1.
func GeneratePUID(msisdn string, length int) (string, error) {
	combinedString := msisdn + uuid.New().String()
	puid := uuid.NewSHA1(uuid.Nil, []byte(combinedString))

	result := strings.ReplaceAll(puid.String(), "-", "")[0:length]
	return result, nil
}
