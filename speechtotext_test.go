package elevenlabs

import (
	"context"
	"testing"
)

func TestTranscriptionRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test empty request
	_, err := client.SpeechToText().Transcribe(ctx, &TranscriptionRequest{})
	if err == nil {
		t.Error("Transcribe() with empty request should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestSpeechToTextService(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Test that service is accessible
	if client.SpeechToText() == nil {
		t.Error("SpeechToText() returned nil")
	}
}

func TestTranscriptionResponse(t *testing.T) {
	// Test response struct initialization
	resp := &TranscriptionResponse{
		Text:         "Hello world",
		LanguageCode: "en",
		Words: []TranscriptionWord{
			{Text: "Hello", Start: 0.0, End: 0.5},
			{Text: "world", Start: 0.6, End: 1.0},
		},
	}

	if resp.Text != "Hello world" {
		t.Errorf("Text = %s, want Hello world", resp.Text)
	}
	if len(resp.Words) != 2 {
		t.Errorf("Words count = %d, want 2", len(resp.Words))
	}
}

// Helper to check if error is ValidationError
func isValidationError(err error, valErr **ValidationError) bool {
	if err == nil {
		return false
	}
	v, ok := err.(*ValidationError)
	if ok && valErr != nil {
		*valErr = v
	}
	return ok
}
