package util

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic string",
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			name:     "string with punctuation",
			input:    "Hello, World! 123.",
			expected: []string{"hello", "world", "123"},
		},
		{
			name:     "ignore single characters",
			input:    "a b c def",
			expected: []string{"def"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tokenize(tt.input)
			if len(result) == 0 && len(tt.expected) == 0 {
				return // both empty, slices are conceptually equal here
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tokenize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateTokens(t *testing.T) {
	fields := []string{"Go Developer", "Backend"}
	tags := []string{"golang", "api"}

	result := GenerateTokens(fields, tags)
	expected := []string{"go", "developer", "backend", "golang", "api"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GenerateTokens() = %v, want %v", result, expected)
	}
}
