package operator

import (
	"testing"
)

func TestTernary(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		a         any
		b         any
		want      any
	}{
		{
			name:      "true returns first value (int)",
			condition: true,
			a:         42,
			b:         0,
			want:      42,
		},
		{
			name:      "false returns second value (int)",
			condition: false,
			a:         42,
			b:         0,
			want:      0,
		},
		{
			name:      "true returns first value (string)",
			condition: true,
			a:         "hello",
			b:         "world",
			want:      "hello",
		},
		{
			name:      "false returns second value (string)",
			condition: false,
			a:         "hello",
			b:         "world",
			want:      "world",
		},
		{
			name:      "true with nil second",
			condition: true,
			a:         "value",
			b:         nil,
			want:      "value",
		},
		{
			name:      "false with nil first",
			condition: false,
			a:         nil,
			b:         "value",
			want:      "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ternary(tt.condition, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Ternary(%v, %v, %v) = %v, want %v", tt.condition, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestTernaryString(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		a         string
		b         string
		want      string
	}{
		{
			name:      "true returns first string",
			condition: true,
			a:         "yes",
			b:         "no",
			want:      "yes",
		},
		{
			name:      "false returns second string",
			condition: false,
			a:         "yes",
			b:         "no",
			want:      "no",
		},
		{
			name:      "true with empty strings",
			condition: true,
			a:         "",
			b:         "fallback",
			want:      "",
		},
		{
			name:      "false with empty strings",
			condition: false,
			a:         "value",
			b:         "",
			want:      "",
		},
		{
			name:      "both same value true",
			condition: true,
			a:         "same",
			b:         "same",
			want:      "same",
		},
		{
			name:      "both same value false",
			condition: false,
			a:         "same",
			b:         "same",
			want:      "same",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TernaryString(tt.condition, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("TernaryString(%v, %q, %q) = %q, want %q", tt.condition, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestTernaryFloat(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		a         float64
		b         float64
		want      float64
	}{
		{
			name:      "true returns first float",
			condition: true,
			a:         3.14,
			b:         2.71,
			want:      3.14,
		},
		{
			name:      "false returns second float",
			condition: false,
			a:         3.14,
			b:         2.71,
			want:      2.71,
		},
		{
			name:      "true with zero",
			condition: true,
			a:         0.0,
			b:         1.0,
			want:      0.0,
		},
		{
			name:      "false with zero",
			condition: false,
			a:         1.0,
			b:         0.0,
			want:      0.0,
		},
		{
			name:      "negative values true",
			condition: true,
			a:         -1.5,
			b:         -2.5,
			want:      -1.5,
		},
		{
			name:      "negative values false",
			condition: false,
			a:         -1.5,
			b:         -2.5,
			want:      -2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TernaryFloat(tt.condition, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("TernaryFloat(%v, %v, %v) = %v, want %v", tt.condition, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestTernary_TypedValues(t *testing.T) {
	t.Run("struct type", func(t *testing.T) {
		type point struct{ X, Y int }
		a := point{1, 2}
		b := point{3, 4}

		got := Ternary(true, a, b)
		if got != a {
			t.Errorf("Ternary(true, %v, %v) = %v, want %v", a, b, got, a)
		}

		got = Ternary(false, a, b)
		if got != b {
			t.Errorf("Ternary(false, %v, %v) = %v, want %v", a, b, got, b)
		}
	})

	t.Run("slice type", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{4, 5, 6}

		got := Ternary(true, a, b)
		if len(got) != 3 || got[0] != 1 {
			t.Errorf("Ternary(true, ...) returned unexpected slice: %v", got)
		}
	})

	t.Run("bool type", func(t *testing.T) {
		got := Ternary(true, true, false)
		if got != true {
			t.Errorf("Ternary(true, true, false) = %v, want true", got)
		}
	})
}
