package operator_test

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/operator"
)

func ExampleTernary() {
	// Generic ternary operator
	result := operator.Ternary(true, "yes", "no")
	fmt.Println(result)

	result2 := operator.Ternary(false, "yes", "no")
	fmt.Println(result2)

	// Works with any type
	num := operator.Ternary(10 > 5, 100, 0)
	fmt.Println(num)

	// Output:
	// yes
	// no
	// 100
}

func ExampleTernaryString() {
	result := operator.TernaryString(true, "active", "inactive")
	fmt.Println(result)

	result2 := operator.TernaryString(false, "active", "inactive")
	fmt.Println(result2)

	// Output:
	// active
	// inactive
}

func ExampleTernaryFloat() {
	result := operator.TernaryFloat(true, 1.5, 2.5)
	fmt.Printf("%.1f\n", result)

	result2 := operator.TernaryFloat(false, 1.5, 2.5)
	fmt.Printf("%.1f\n", result2)

	// Output:
	// 1.5
	// 2.5
}
