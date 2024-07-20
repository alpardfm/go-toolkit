package validation

import (
	"regexp"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

func IsValidEmail(email string) (bool, error) {
	const pattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	isMatch, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false, err
	}

	if !isMatch {
		return false, errors.NewWithCode(codes.CodeInvalidValue, "email format is not valid")
	} else {
		return true, nil
	}
}
