package encrypt

import (
	"encoding/base64"
	"testing"
)

func TestAES256GCM(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("this is key string"))
	tests := []struct {
		name       string
		text       string
		wantResult string
		isEncrypt  bool
		isDecrypt  bool
	}{
		{
			name:      "test encrypt",
			text:      "password",
			isEncrypt: true,
		},
		{
			name:       "test decrypt",
			text:       "MWM0Y2UwNDQzMzMwYmUzMzYzNzIyMDZjZjIyNTUwOWVhMmYxZTJkZjg2NDk3OGFjNzdmMGU4YWI3ZjdiMjhkNzhmYzA4MWEy",
			wantResult: "password",
			isDecrypt:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isDecrypt {
				decryptedText, err := DecryptAES256GC([]byte(test.text), []byte(key))
				if err != nil {
					t.Errorf("TestAES256GCM() error: '%v', want err: 'nil'", err)
				} else {
					if string(decryptedText) != test.wantResult {
						t.Errorf("TestAES256GCM() result: '%v', want result: '%v'", string(decryptedText), test.wantResult)
					}
				}
			}

			if test.isEncrypt {
				encryptedText, err := EncryptAES256GCM([]byte(test.text), []byte(key))
				if err != nil {
					t.Errorf("TestAES256GCM() error: '%v', want err: 'nil'", err)
				} else {
					if string(encryptedText) == "" {
						t.Errorf("TestAES256GCM() error: '%v', want err: 'nil'", "encrypted text is empty")
					} else if string(encryptedText) == test.text {
						t.Errorf("TestAES256GCM() error: '%v', want err: '{encrypted text}'", "encrypted text is equals with plain text")
					}
				}
			}
		})
	}

}
