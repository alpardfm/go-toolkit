package hash

import (
	"strings"
	"testing"
)

func TestHashBcrypt(t *testing.T) {
	defaultPassword := "password"
	tests := []struct {
		name       string
		password   string
		cost       int
		messageErr string
		wantErr    bool
		mustEquals bool
	}{
		{
			name:       "test hash bcrypt",
			password:   defaultPassword,
			cost:       10,
			messageErr: "",
			wantErr:    false,
			mustEquals: true,
		},
		{
			name:       "test invalid hash bcrypt because cost greather than max cost",
			password:   defaultPassword,
			cost:       35,
			messageErr: "error while encode bcrypt hash, crypto/bcrypt: cost 35 is outside allowed range (4,31)",
			wantErr:    true,
			mustEquals: false,
		},
		{
			name:       "test invalid hash bcrypt because cost lower than min cost",
			password:   defaultPassword,
			cost:       3,
			messageErr: "error while encode bcrypt hash, crypto/bcrypt: cost 3 is outside allowed range (4,31)",
			wantErr:    true,
			mustEquals: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := CreateBcrypt(test.password, test.cost)
			if err != nil {
				if test.wantErr {
					if !strings.Contains(test.messageErr, err.Error()) {
						t.Errorf(err.Error())
					}
				}
			}
		})
	}
}

func TestCompareBcrypt(t *testing.T) {
	defaultPassword := "1234567890"
	tests := []struct {
		name         string
		hashPassword string
		password     string
		messageErr   string
		wantErr      bool
		mustEquals   bool
	}{
		{
			name:         "test match when compare bcrypt",
			hashPassword: "$2b$10$8ImmbtD6jb4J6SxU2SIAjeSgJvABYwsKRhMRR4pRQHMekqTWR3Mh.",
			password:     defaultPassword,
			messageErr:   "",
			wantErr:      false,
			mustEquals:   true,
		},
		{
			name:         "test invalid password when compare bcrypt",
			hashPassword: "$2b$10$8ImmbtD6jb4J6SxU2SIAjeSgJvABYwsKRhMRR4pRQHMekqTWR3Mh.",
			password:     "123456789",
			messageErr:   "error while compare bcrypt hash, crypto/bcrypt: hashedPassword is not the hash of the given password",
			wantErr:      true,
			mustEquals:   false,
		},
		{
			name:         "test invalid hash password when compare bcrypt (lt)",
			hashPassword: "$2b$10$8ImmbtD6jb4J6SxU2SIAjeSgJvABYwsKRhMRR4pRQHMekq",
			password:     defaultPassword,
			messageErr:   "error while compare bcrypt hash, crypto/bcrypt: hashedSecret too short to be a bcrypted password",
			wantErr:      true,
			mustEquals:   false,
		},
		{
			name:         "test invalid hash password when compare bcrypt (wrong)",
			hashPassword: "$2b$10$8ImmbtD6jb4J6SxU2SIAjeSgJvABYwsKRhMRR4pRQHMekqTWR3Mh,",
			password:     defaultPassword,
			messageErr:   "error while compare bcrypt hash, crypto/bcrypt: hashedPassword is not the hash of the given password",
			wantErr:      true,
			mustEquals:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CompareBcrypt(test.hashPassword, test.password)
			if err != nil {
				if test.wantErr {
					if !strings.Contains(test.messageErr, err.Error()) {
						t.Errorf(err.Error())
					}
				}
			}
		})
	}
}
