package elevenlabs

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "voice_id", Message: "cannot be empty"}
	expected := "elevenlabs: validation error for voice_id: cannot be empty"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %s, want %s", err.Error(), expected)
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "without detail",
			err:      &APIError{StatusCode: 401, Message: "Unauthorized"},
			expected: "elevenlabs: API error (status 401): Unauthorized",
		},
		{
			name:     "with detail",
			err:      &APIError{StatusCode: 400, Message: "Bad Request", Detail: "Invalid voice_id"},
			expected: "elevenlabs: API error (status 400): Bad Request - Invalid voice_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("APIError.Error() = %s, want %s", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "404 error",
			err:      &APIError{StatusCode: 404, Message: "Not Found"},
			expected: true,
		},
		{
			name:     "401 error",
			err:      &APIError{StatusCode: 401, Message: "Unauthorized"},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.err); got != tt.expected {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "401 error",
			err:      &APIError{StatusCode: 401, Message: "Unauthorized"},
			expected: true,
		},
		{
			name:     "404 error",
			err:      &APIError{StatusCode: 404, Message: "Not Found"},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorizedError(tt.err); got != tt.expected {
				t.Errorf("IsUnauthorizedError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "429 error",
			err:      &APIError{StatusCode: 429, Message: "Too Many Requests"},
			expected: true,
		},
		{
			name:     "401 error",
			err:      &APIError{StatusCode: 401, Message: "Unauthorized"},
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimitError(tt.err); got != tt.expected {
				t.Errorf("IsRateLimitError() = %v, want %v", got, tt.expected)
			}
		})
	}
}
