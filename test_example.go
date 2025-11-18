package main

import "fmt"

// This is a test file with some issues for the AI to catch
func calculateSum(numbers []int) int {
	sum := 0
	for i := 0; i < len(numbers); i++ {
		sum += numbers[i]
	}
	return sum
}

// Function with unused parameter
func greet(name string, age int) string {
	return "Hello, " + name
}

// Missing error handling
func divide(a, b int) int {
	return a / b
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	total := calculateSum(nums)
	fmt.Println("Sum:", total)

	// Potential panic - division by zero
	result := divide(10, 0)
	fmt.Println("Result:", result)
}
