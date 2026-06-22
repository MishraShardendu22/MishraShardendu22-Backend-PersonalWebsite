package util

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	val := GetEnv("TEST_KEY", "fallback")
	if val != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", val)
	}

	// Test with missing environment variable
	val = GetEnv("MISSING_KEY", "fallback")
	if val != "fallback" {
		t.Errorf("Expected 'fallback', got '%s'", val)
	}
}
