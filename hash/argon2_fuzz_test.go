package hash

import (
	"testing"
)

// FuzzDecodeHash tests that decodeHash never panics on arbitrary input.
func FuzzDecodeHash(f *testing.F) {
	// Seed corpus with valid and edge-case inputs
	f.Add("$argon2id$v=19$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g")
	f.Add("")
	f.Add("$$$$$")
	f.Add("$argon2id$v=19$m=4096,t=3,p=1$invalid-base64!$also-invalid!")
	f.Add("$argon2id$v=99$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g")
	f.Add("not-a-hash-at-all")
	f.Add("$argon2id$v=abc$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g")

	f.Fuzz(func(t *testing.T, input string) {
		// Should never panic, errors are acceptable
		_, _, _ = decodeHash(input)
	})
}

// FuzzCompareArgon2 tests that CompareArgon2 never panics on arbitrary input.
func FuzzCompareArgon2(f *testing.F) {
	f.Add("password", "$argon2id$v=19$m=4096,t=3,p=1$c29tZXNhbHQ$c29tZWhhc2g")
	f.Add("", "")
	f.Add("test", "garbage")

	f.Fuzz(func(t *testing.T, password, encodedHash string) {
		// Should never panic
		_, _ = CompareArgon2(password, encodedHash)
	})
}
