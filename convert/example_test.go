package convert_test

import (
	"fmt"

	"github.com/alpardfm/go-toolkit/convert"
)

func ExampleToInt64() {
	val, err := convert.ToInt64("42")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(val)

	// Output:
	// 42
}

func ExampleToString() {
	val, err := convert.ToString(123)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(val)

	// Output:
	// 123
}

func ExampleToPtr() {
	s := "hello"
	ptr := convert.ToPtr(s)
	fmt.Println(*ptr)

	num := 42
	numPtr := convert.ToPtr(num)
	fmt.Println(*numPtr)

	// Output:
	// hello
	// 42
}

func ExampleToFloat64() {
	val, err := convert.ToFloat64("3.14")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%.2f\n", val)

	// Output:
	// 3.14
}

func ExampleIntToChar() {
	fmt.Println(convert.IntToChar(1))
	fmt.Println(convert.IntToChar(3))
	fmt.Println(convert.IntToChar(27))

	// Output:
	// A
	// C
	// AA
}

func ExampleToSafeValue() {
	// Returns zero value for nil pointer
	var nilPtr *string
	result := convert.ToSafeValue[string](nilPtr)
	fmt.Printf("%q\n", result)

	// Returns actual value for valid input
	val := "hello"
	result2 := convert.ToSafeValue[string](&val)
	fmt.Println(result2)

	// Output:
	// ""
	// hello
}
