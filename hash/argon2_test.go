package hash

import (
	"strings"
	"testing"
)

func TestHashAndCompareArgon2(t *testing.T) {
	defaultPassword := "password"
	tests := []struct {
		name               string
		password           string
		wantErr            bool
		wantErrMsgContains string
		wantResultContains string
		mustEquals         bool
	}{
		{
			name:               "test hash and compare password equals",
			password:           defaultPassword,
			wantErr:            false,
			wantErrMsgContains: "",
			wantResultContains: "argon2id$v=",
			mustEquals:         true,
		},
		{
			name:               "test hash and compare password not equals",
			password:           defaultPassword,
			wantErr:            true,
			wantErrMsgContains: "",
			wantResultContains: "argon2id$v=",
			mustEquals:         false,
		},
		{
			name:               "test invalid length password",
			password:           defaultPassword,
			wantErr:            true,
			wantErrMsgContains: "invalid length of encoded hash",
			wantResultContains: "argon2id$v=",
			mustEquals:         false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// try to hash
			hash, err := NewArgon2([]byte(test.password))
			if err != nil {
				// if want error
				if test.wantErr {
					if !strings.Contains(err.Error(), test.wantErrMsgContains) {
						t.Errorf("NewArgon2() error: '%v', want error: '%v'", err, test.wantErrMsgContains)
					}
				} else {
					t.Errorf("NewArgon2() error: '%v' want error: nil", err)
				}
			} else {
				if !strings.Contains(hash, test.wantResultContains) {
					t.Errorf("NewArgon2() want contains: '%s', get: '%s'", test.wantResultContains, hash)
				} else {
					// if must not equals
					if !test.mustEquals {
						test.password = strings.Repeat("x", len(test.password))
					}
					// if must equals
					isEquals, err := CompareArgon2(test.password, hash)
					if err != nil && !strings.Contains(err.Error(), test.wantErrMsgContains) {
						t.Errorf("NewArgon2() error: '%v', want err: nil", err)
					} else {
						if test.mustEquals {
							if !isEquals {
								t.Errorf("NewArgon2() result: %v, want result: %v", isEquals, test.mustEquals)
							}
						} else {
							if isEquals {
								t.Errorf("NewArgon2() result: %v, want result: %v", isEquals, test.mustEquals)
							}
						}
					}

				}

			}
		})
	}
}
