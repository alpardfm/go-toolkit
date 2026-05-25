package random_test

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/random"
)

func ExampleGenerateInt() {
	// Generate a random integer with 6 digits
	val, err := random.GenerateInt(6)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// The result is a positive integer with exactly 6 digits
	fmt.Println(val >= 100000 && val <= 999999)

	// Output:
	// true
}
