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

func GenerateInt(length int) (int64, error) {
	rand.Seed(time.Now().UnixNano())
	value := ""
	for i := 0; i < length; i++ {
		s, err := convert.ToString(rand.Intn(9))
		if err != nil {
			return 0, errors.NewWithCode(codes.CodeInvalidValue, err.Error())
		}

		if s == "0" {
			i -= 1
			continue
		}

		value += s
	}

	iValue, err := convert.ToInt64(value)
	if err != nil {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, err.Error())
	}

	return iValue, nil
}

func GeneratePUID(msisdn string, length int) (string, error) {
	combinedString := msisdn + uuid.New().String()
	puid := uuid.NewSHA1(uuid.Nil, []byte(combinedString))

	result := strings.ReplaceAll(puid.String(), "-", "")[0:length]
	return result, nil
}
