package elevenlabs

import (
	"context"
	"testing"
)

func TestDialogueRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test empty inputs
	_, err := client.TextToDialogue().Generate(ctx, &DialogueRequest{})
	if err == nil {
		t.Error("Generate() with empty inputs should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "inputs" {
		t.Errorf("ValidationError field = %s, want inputs", valErr.Field)
	}
}

func TestTextToDialogueService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.TextToDialogue() == nil {
		t.Error("TextToDialogue() returned nil")
	}
}

func TestDialogueInput(t *testing.T) {
	// Test DialogueInput struct
	input := DialogueInput{
		VoiceID: "voice123",
		Text:    "Hello world",
	}

	if input.VoiceID != "voice123" {
		t.Errorf("VoiceID = %s, want voice123", input.VoiceID)
	}
	if input.Text != "Hello world" {
		t.Errorf("Text = %s, want Hello world", input.Text)
	}
}

func TestDialogueResponse(t *testing.T) {
	// Test response struct initialization
	resp := &DialogueResponse{
		AudioBase64: "base64data",
		VoiceSegments: []VoiceSegment{
			{VoiceID: "voice1", StartTime: 0.0, EndTime: 1.5},
			{VoiceID: "voice2", StartTime: 1.6, EndTime: 3.0},
		},
	}

	if resp.AudioBase64 != "base64data" {
		t.Errorf("AudioBase64 = %s, want base64data", resp.AudioBase64)
	}
	if len(resp.VoiceSegments) != 2 {
		t.Errorf("VoiceSegments count = %d, want 2", len(resp.VoiceSegments))
	}
	if resp.VoiceSegments[0].VoiceID != "voice1" {
		t.Errorf("VoiceSegments[0].VoiceID = %s, want voice1", resp.VoiceSegments[0].VoiceID)
	}
}

func TestSimpleDialogue(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test Simple with empty inputs
	_, err := client.TextToDialogue().Simple(ctx, nil)
	if err == nil {
		t.Error("Simple() with nil inputs should return error")
	}

	_, err = client.TextToDialogue().Simple(ctx, []DialogueInput{})
	if err == nil {
		t.Error("Simple() with empty inputs should return error")
	}
}

func TestGenerateStream(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test GenerateStream with empty inputs
	_, err := client.TextToDialogue().GenerateStream(ctx, &DialogueRequest{})
	if err == nil {
		t.Error("GenerateStream() with empty inputs should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}
