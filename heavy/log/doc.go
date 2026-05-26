// Package log provides structured logging via zerolog with automatic context
// field extraction (request ID, user agent, service version, elapsed time).
//
// Each Init() call returns a new independent logger instance, making it safe
// to use in tests and to create multiple loggers with different configurations.
// Use the Config.Writer field to redirect output (defaults to os.Stdout).
package log
