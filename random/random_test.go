package random

import (
	"fmt"
	"testing"
)

func TestGenerateInt(t *testing.T) {
	value := ""
	for i := 1; i <= 1000; i++ {
		value = func() string {
			v, err := GenerateInt(6)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			return fmt.Sprintf("%v", v)
		}()

		if len(value) < 6 {
			t.Fatalf("error: length GenRandomInt is lower than 6")
		} else {
			t.Logf("success: GenRandomInt result length is 6 : %v", value)
		}
	}
}

func TestGeneratePuid(t *testing.T) {
	tests := []struct {
		msisdn string
		length int
	}{
		{
			msisdn: "087864235314",
			length: 14,
		},
		{
			msisdn: "087864235314",
			length: 13,
		},
		{
			msisdn: "085713264497",
			length: 15,
		},
	}

	uuid := map[string]string{}

	for _, test := range tests {
		t.Run("Generate UUID with Misdn : "+test.msisdn, func(t *testing.T) {
			// Todo generate puid
			// Todo check is string length
			// check in map is key with generated PUID is already exists? t.Errorf if exists and don't do anything if doesn't exists

			generate, err := GeneratePUID(test.msisdn, test.length)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			if len(generate) != test.length {
				t.Errorf("Length String Not Equals")
			}

			_, exits := uuid[generate]

			if exits {
				t.Errorf("duplicate uuid : " + generate)
			} else {
				uuid[generate] = "ashdkashdkjhsd"
				t.Logf("Generate : %v", generate)
			}
		})
	}
}
