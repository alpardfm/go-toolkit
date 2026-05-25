package errors_test

import (
	"fmt"
	"io"
	"strings"

	stderrors "errors"

	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
)

func ExampleNewWithCode() {
	// Create a coded error
	err := errors.NewWithCode(codes.CodeInvalidValue, "invalid input: %s", "email is required")
	fmt.Println(err)

	// Get the code from the error
	code := errors.GetCode(err)
	fmt.Println(code)

	// Output:
	// Error: invalid input: email is required
	// 1000
}

func ExampleGetCode() {
	err := errors.NewWithCode(codes.CodeBadRequest, "bad request")
	code := errors.GetCode(err)
	fmt.Println(code)

	// Non-toolkit errors return NoCode (MaxUint32)
	plainErr := fmt.Errorf("plain error")
	code2 := errors.GetCode(plainErr)
	fmt.Println(code2)

	// Output:
	// 1006
	// 4294967295
}

func ExampleCompile() {
	err := errors.NewWithCode(codes.CodeAuthFailure, "token expired")
	statusCode, appErr := errors.Compile(err, "en")
	fmt.Println(statusCode)
	fmt.Println(appErr.Code)
	fmt.Println(appErr.Title)

	// Output:
	// 401
	// 1703
	// Unauthorized
}

// ExampleNewWithCode_unwrapChain demonstrates creating an error chain
// and traversing it with standard library errors.Is.
func ExampleNewWithCode_unwrapChain() {
	// Create a sentinel error
	sentinel := io.EOF

	// Demonstrate that GetCode works on toolkit errors
	err1 := errors.NewWithCode(codes.CodeBadRequest, "first error")
	err2 := errors.NewWithCode(codes.CodeInvalidValue, "second error")
	fmt.Println(errors.GetCode(err1))
	fmt.Println(errors.GetCode(err2))

	// Standard errors.Is works with sentinel errors
	wrappedSentinel := fmt.Errorf("wrapped: %w", sentinel)
	fmt.Println(stderrors.Is(wrappedSentinel, io.EOF))

	// Verify the error message contains expected text
	fmt.Println(strings.Contains(err1.Error(), "first error"))

	// Output:
	// 1006
	// 1000
	// true
	// true
}
