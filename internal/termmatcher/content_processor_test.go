package termmatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultContentProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic processing",
			input:    "The quick brown fox jumps over the lazy dog",
			expected: "quick brown fox jump lazi dog",
		},
		{
			name:     "With hyphens",
			input:    "state-of-the-art technology",
			expected: "state art technolog",
		},
		{
			name:     "With stopwords",
			input:    "The cat is on the mat",
			expected: "cat mat",
		},
		{
			name:     "Mixed case",
			input:    "ThE QuIcK BrOwN FoX",
			expected: "quick brown fox",
		},
	}

	processor := NewDefaultContentProcessor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultContentProcessor_Stem(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic stemming",
			input:    "running jumps foxes",
			expected: "run jump fox",
		},
		{
			name:     "Mixed case",
			input:    "Jumping RUNNING Foxes",
			expected: "jump run fox",
		},
	}

	processor := NewDefaultContentProcessor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Stem(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultContentProcessor_ProcessCaching(t *testing.T) {
	processor := NewDefaultContentProcessor()

	input := "The quick brown fox"
	expected := "quick brown fox"

	// First call should process the content
	result1 := processor.Process(input)
	assert.Equal(t, expected, result1)

	// Second call should return the cached result
	result2 := processor.Process(input)
	assert.Equal(t, expected, result2)

	// Verify that the results are the same
	assert.Equal(t, result1, result2)
}
