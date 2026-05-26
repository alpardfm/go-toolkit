package validation

import (
	"regexp"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

// emailRegex is compiled once at package init for performance.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail validates whether the given string is a valid email address format.
// Returns true if valid, or false with an error if the format is invalid.
func IsValidEmail(email string) (bool, error) {
	if !emailRegex.MatchString(email) {
		return false, errors.NewWithCode(codes.CodeInvalidValue, "email format is not valid")
	}
	return true, nil
}
