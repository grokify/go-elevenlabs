package elevenlabs

import (
	"context"
	"strings"
	"testing"
)

func TestAudioIsolationRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test nil audio
	_, err := client.AudioIsolation().Isolate(ctx, &AudioIsolationRequest{})
	if err == nil {
		t.Error("Isolate() with nil audio should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "audio" {
		t.Errorf("ValidationError field = %s, want audio", valErr.Field)
	}
}

func TestAudioIsolationService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.AudioIsolation() == nil {
		t.Error("AudioIsolation() returned nil")
	}
}

func TestIsolateFile(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test with valid parameters (will fail at API but pass validation)
	audio := strings.NewReader("fake audio data")
	_, err := client.AudioIsolation().IsolateFile(ctx, audio, "test.mp3")
	// This will fail because we don't have a real API key or valid audio
	// but it tests that the method exists and accepts the right parameters
	if err == nil {
		// If no error, the service would be making a real API call
		t.Log("IsolateFile() called successfully")
	}
}

func TestIsolateStream(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test nil audio for stream
	_, err := client.AudioIsolation().IsolateStream(ctx, &AudioIsolationRequest{})
	if err == nil {
		t.Error("IsolateStream() with nil audio should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}
