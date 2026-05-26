// Package codes defines a registry of numeric error and status codes used
// throughout the toolkit. Each code maps to a bilingual (English/Indonesian)
// user-facing message with an associated HTTP status code.
//
// Application code should use these codes with errors.NewWithCode() to create
// errors that can be compiled into user-friendly API responses via errors.Compile().
package codes
