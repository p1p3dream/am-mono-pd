package stringutil

import (
	"strings"
)

// ASCIIToSnakeCase converts a string to snake_case.
// Only handles ASCII inputs.
func ASCIIToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	var previousWasUnderscore bool

	for i := range len(s) {
		char := s[i]

		if char > 127 {
			continue
		}

		// Check if the character is a letter or digit.
		isAlphanumeric := (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')

		if !isAlphanumeric {
			if previousWasUnderscore {
				// Do not add multiple underscores.
				continue
			}

			// Replace non-alphanumeric with underscore
			// if we're not writing the first character.
			if result.Len() > 0 {
				result.WriteByte('_')
				previousWasUnderscore = true
			}

			continue
		}

		if char < 'A' || char > 'Z' {
			// Lowercase char or number, since we're ignoring non-alphanumeric.
			result.WriteByte(char)
			previousWasUnderscore = false
			continue
		}

		// Add underscore before capital letters, but not at the beginning.
		if result.Len() > 0 && !previousWasUnderscore {
			// Check if previous character was lowercase or if we're entering an acronym.
			isPreviousLower := i > 0 && s[i-1] >= 'a' && s[i-1] <= 'z'
			isNextLower := i < len(s)-1 && s[i+1] >= 'a' && s[i+1] <= 'z'

			// Add underscore if previous char was lowercase or if we're at the end of an acronym.
			if isPreviousLower || (i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z' && isNextLower) {
				result.WriteByte('_')
			}
		}

		// Convert to lowercase.
		// ASCII 'A' to 'a' offset is 32.
		result.WriteByte(char + 32)
		previousWasUnderscore = false
	}

	return result.String()
}

func JoinNonEmpty(sep string, parts ...string) string {
	var b strings.Builder

	for _, part := range parts {
		if part == "" {
			continue
		}

		if b.Len() > 0 {
			b.WriteString(sep)
		}

		b.WriteString(part)
	}

	return b.String()
}
