package operator

// TernaryString returns the first value if the condition is true, and returns the second value if the condition is false
func TernaryString(condition bool, a, b string) (str string) {
	if condition {
		str = a
	} else {
		str = b
	}

	return str
}

// TernaryFloat returns the first value if the condition is true, and returns the second value if the condition is false
func TernaryFloat(condition bool, a, b float64) (res float64) {
	if condition {
		res = a
	} else {
		res = b
	}

	return res
}

// Ternary returns the first value if the condition is true, and returns the second value if the the condition is false
func Ternary[T any](condition bool, a, b T) (res T) {
	if condition {
		res = a
	} else {
		res = b
	}

	return res
}
