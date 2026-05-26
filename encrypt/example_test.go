package encrypt_test

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/encrypt"
)

func ExampleEncryptAES256GCM() {
	// AES-256 requires a 32-byte key
	key := []byte("01234567890123456789012345678901")
	plaintext := []byte("hello, world!")

	// Encrypt the plaintext
	ciphertext, err := encrypt.EncryptAES256GCM(plaintext, key)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Decrypt it back
	decrypted, err := encrypt.DecryptAES256GCM(ciphertext, key)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(decrypted))

	// Output:
	// hello, world!
}
