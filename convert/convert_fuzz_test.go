package convert

import (
	"testing"
)

// FuzzToInt64 tests that ToInt64 never panics on arbitrary string input.
func FuzzToInt64(f *testing.F) {
	f.Add("123")
	f.Add("-456")
	f.Add("0")
	f.Add("")
	f.Add("not-a-number")
	f.Add("9999999999999999999")
	f.Add("3.14")
	f.Add("1e10")

	f.Fuzz(func(t *testing.T, input string) {
		// Should never panic
		_, _ = ToInt64(input)
	})
}

// FuzzToFloat64 tests that ToFloat64 never panics on arbitrary string input.
func FuzzToFloat64(f *testing.F) {
	f.Add("3.14")
	f.Add("-2.5")
	f.Add("0")
	f.Add("")
	f.Add("not-a-number")
	f.Add("1e308")
	f.Add("inf")
	f.Add("NaN")

	f.Fuzz(func(t *testing.T, input string) {
		// Should never panic
		_, _ = ToFloat64(input)
	})
}

// FuzzToString tests that ToString never panics on arbitrary int input.
func FuzzToString(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(123))
	f.Add(int64(-456))
	f.Add(int64(9223372036854775807))  // max int64
	f.Add(int64(-9223372036854775808)) // min int64

	f.Fuzz(func(t *testing.T, input int64) {
		// Should never panic
		result, err := ToString(input)
		if err != nil {
			t.Fatalf("ToString(int64) should never fail: %v", err)
		}
		if result == "" {
			t.Fatal("ToString(int64) returned empty string")
		}
	})
}

// FuzzIntToChar tests that IntToChar never panics on arbitrary input.
func FuzzIntToChar(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(26)
	f.Add(27)
	f.Add(702)
	f.Add(703)
	f.Add(-1)
	f.Add(10000)

	f.Fuzz(func(t *testing.T, input int) {
		// Should never panic (even for negative or very large values)
		_ = IntToChar(input)
	})
}
