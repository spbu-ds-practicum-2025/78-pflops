package utils

import "strings"

func IsValidEmail(email string) bool {
	// Simple email validation logic
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func IsValidPassword(password string) bool {
	// Simple password validation logic
	return len(password) >= 8 && strings.ContainsAny(password, "!@#$%^&*()_+")
}
