package utils

import "strings"

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference of two integers
func Subtract(a, b int) int {
	return a - b
}

// Divide returns the quotient of two integers
// If b is zero, it returns zero to avoid panic
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

// ReverseString returns the reverse of a given string
func ReverseString(s string) string {
	var reversed strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		reversed.WriteByte(s[i])
	}
	return reversed.String()
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
	reversed := ReverseString(s)
	return s == reversed
}
