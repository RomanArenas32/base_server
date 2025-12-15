package pkg

// Add adds two integers and returns the result.
func Add(a, b int) int {
	return a + b
}

// Subtract subtracts b from a and returns the result.
func Subtract(a, b int) int {
	return a - b
}

// Multiply multiplies two integers and returns the result.
func Multiply(a, b int) int {
	return a * b
}

// Divide divides a by b and returns the result.
// Returns 0 if b is 0 to avoid division by zero.
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
