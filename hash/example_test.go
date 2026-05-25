package hash_test

import (
	"fmt"
	"strings"

	"github.com/alpardfm/go-toolkit/hash"
)

func ExampleNewArgon2() {
	// Hash a password using Argon2id
	hashed, err := hash.NewArgon2([]byte("my-secret-password"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// The hash follows the standard argon2id format
	fmt.Println(strings.HasPrefix(hashed, "$argon2id$"))

	// Output:
	// true
}

func ExampleCompareArgon2() {
	// First, hash a password
	password := "my-secret-password"
	hashed, err := hash.NewArgon2([]byte(password))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Compare with the correct password
	match, err := hash.CompareArgon2(password, hashed)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("correct password:", match)

	// Compare with a wrong password
	match, err = hash.CompareArgon2("wrong-password", hashed)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("wrong password:", match)

	// Output:
	// correct password: true
	// wrong password: false
}
