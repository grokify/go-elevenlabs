package elevenlabs

import (
	"context"
	"testing"
)

func TestMusicRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test empty prompt
	_, err := client.Music().Generate(ctx, &MusicRequest{})
	if err == nil {
		t.Error("Generate() with empty prompt should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "prompt" {
		t.Errorf("ValidationError field = %s, want prompt", valErr.Field)
	}
}

func TestMusicService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.Music() == nil {
		t.Error("Music() returned nil")
	}
}

func TestMusicRequest(t *testing.T) {
	// Test MusicRequest struct
	req := MusicRequest{
		Prompt:            "upbeat electronic music",
		DurationMs:        30000,
		ForceInstrumental: true,
		Seed:              12345,
	}

	if req.Prompt != "upbeat electronic music" {
		t.Errorf("Prompt = %s, want upbeat electronic music", req.Prompt)
	}
	if req.DurationMs != 30000 {
		t.Errorf("DurationMs = %d, want 30000", req.DurationMs)
	}
	if !req.ForceInstrumental {
		t.Error("ForceInstrumental should be true")
	}
	if req.Seed != 12345 {
		t.Errorf("Seed = %d, want 12345", req.Seed)
	}
}

func TestMusicResponse(t *testing.T) {
	// Test MusicResponse struct
	resp := MusicResponse{
		SongID: "song123",
	}

	if resp.SongID != "song123" {
		t.Errorf("SongID = %s, want song123", resp.SongID)
	}
}

func TestGenerateStreamValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test GenerateStream with empty prompt
	_, err := client.Music().GenerateStream(ctx, &MusicRequest{})
	if err == nil {
		t.Error("GenerateStream() with empty prompt should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "prompt" {
		t.Errorf("ValidationError field = %s, want prompt", valErr.Field)
	}
}

func TestSimpleMusicValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test Simple with empty prompt
	_, err := client.Music().Simple(ctx, "")
	if err == nil {
		t.Error("Simple() with empty prompt should return error")
	}
}

func TestGenerateInstrumentalValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test GenerateInstrumental with empty prompt
	_, err := client.Music().GenerateInstrumental(ctx, "", 30000)
	if err == nil {
		t.Error("GenerateInstrumental() with empty prompt should return error")
	}
}
