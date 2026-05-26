// Package errors provides stacktrace-enriched error handling with error codes.
//
// Errors created with NewWithCode capture the file, line number, and function name
// at the call site. They carry a numeric code (from the codes package) that maps
// to user-facing messages in multiple languages.
//
// Use GetCaller to extract stacktrace information, and GetCode to retrieve the
// error code from any error in the chain. Use Compile to convert an error into
// a user-facing response with the appropriate HTTP status code and localized message.
package errors
