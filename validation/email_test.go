package validation

import (
	"testing"

	"github.com/alpardfm/go-toolkit/convert"
)

func TestIsValidEmail(t *testing.T) {
	method := "IsValidEmail()"
	tests := []struct {
		name        string
		paramsEmail string
		wantErrMsg  string
		wantResult  *bool
	}{
		{
			name:        "test email valid",
			paramsEmail: "example@gmail.com",
			wantErrMsg:  "",
			wantResult:  convert.ToPtr(true),
		},
		{
			name:        "test email invalid format",
			paramsEmail: "examplegmail.com",
			wantErrMsg:  "email format is not valid",
			wantResult:  convert.ToPtr(false),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isMatch, err := IsValidEmail(test.paramsEmail)
			if err != nil {
				if test.wantErrMsg != "" {
					if err.Error() != test.wantErrMsg {
						t.Errorf("%v wantErrMsg: '%v', getErrMsg: '%v'", method, test.wantErrMsg, err)
					}
				} else {
					t.Errorf("%v wantErrMsg: '%v', getErrMsg: '%v'", method, test.wantErrMsg, err)
				}
			}

			if test.wantResult != nil {
				if isMatch != *test.wantResult {
					t.Errorf("%v wantResult: '%v', getResult: '%v'", method, *test.wantResult, isMatch)
				}
			}
		})
	}
}
