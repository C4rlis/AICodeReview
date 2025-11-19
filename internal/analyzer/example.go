package analyzer

import (
	"fmt"
)

// Example function with some deliberate issues for AI to review
func ProcessData(data []string) ([]string, error) {
	// TODO: Add input validation
	result := make([]string, 0)

	for i := 0; i < len(data); i++ {
		// Non-idiomatic loop
		if data[i] != "" {
			result = append(result, data[i])
		}
	}

	return result, nil
}

// Missing context usage
func FetchUserData(userID int) (string, error) {
	// Should use context for cancellation
	fmt.Printf("Fetching user %d\n", userID)
	return fmt.Sprintf("User-%d", userID), nil
}

// Function with unused parameter
func ValidateInput(input string, maxLength int) bool {
	// maxLength is not used!
	return len(input) > 0
}
