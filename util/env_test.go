package util

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	originalValue := os.Getenv("TEST_KEY")

	// Set a test value for the environment variable
	testValue := "test_value"
	os.Setenv("TEST_KEY", testValue)

	// Call the Env function and check if it returns the expected value
	result := Env("TEST_KEY")
	if result != testValue {
		t.Errorf("Expected %s, got %s", testValue, result)
	}

	// Unset the environment variable to clean up
	os.Unsetenv("TEST_KEY")

	// Call the Env function again, expecting an empty string since the environment variable is unset
	result = Env("TEST_KEY")
	if result != "" {
		t.Errorf("Expected an empty string, got %s", result)
	}

	// Restore the original value of the environment variable
	os.Setenv("TEST_KEY", originalValue)
}
