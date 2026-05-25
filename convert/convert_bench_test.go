package convert

import (
	"testing"
)

func BenchmarkToInt64(b *testing.B) {
	b.ReportAllocs()
	input := "1234567890"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ToInt64(input)
	}
}

func BenchmarkToFloat64(b *testing.B) {
	b.ReportAllocs()
	input := "3.141592653589793"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ToFloat64(input)
	}
}

func BenchmarkToString(b *testing.B) {
	b.ReportAllocs()
	input := int64(1234567890)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ToString(input)
	}
}

func BenchmarkToSafeValue(b *testing.B) {
	b.ReportAllocs()
	input := "hello world"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ToSafeValue[string](input)
	}
}
