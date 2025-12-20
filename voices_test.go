package elevenlabs

import (
	"context"
	"testing"
)

func TestVoicesList_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	voices, err := client.Voices().List(context.Background())
	if err != nil {
		t.Fatalf("Voices().List() error = %v", err)
	}
	if len(voices) == 0 {
		t.Error("Voices().List() returned empty list")
	}

	// Check that voices have required fields
	for _, v := range voices {
		if v.VoiceID == "" {
			t.Error("Voice has empty VoiceID")
		}
		if v.Name == "" {
			t.Error("Voice has empty Name")
		}
	}
}

func TestVoicesGet_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// First get list to find a valid voice ID
	voices, err := client.Voices().List(context.Background())
	if err != nil {
		t.Fatalf("Voices().List() error = %v", err)
	}
	if len(voices) == 0 {
		t.Skip("No voices available")
	}

	voiceID := voices[0].VoiceID

	voice, err := client.Voices().Get(context.Background(), voiceID)
	if err != nil {
		t.Fatalf("Voices().Get() error = %v", err)
	}
	if voice == nil {
		t.Fatal("Voices().Get() returned nil")
	}
	if voice.VoiceID != voiceID {
		t.Errorf("Voice.VoiceID = %s, want %s", voice.VoiceID, voiceID)
	}
}

func TestVoicesGetSettings_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// First get list to find a valid voice ID
	voices, err := client.Voices().List(context.Background())
	if err != nil {
		t.Fatalf("Voices().List() error = %v", err)
	}
	if len(voices) == 0 {
		t.Skip("No voices available")
	}

	voiceID := voices[0].VoiceID

	settings, err := client.Voices().GetSettings(context.Background(), voiceID)
	if err != nil {
		t.Fatalf("Voices().GetSettings() error = %v", err)
	}
	if settings == nil {
		t.Fatal("Voices().GetSettings() returned nil")
	}
}

func TestVoicesGetDefaultSettings_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	settings, err := client.Voices().GetDefaultSettings(context.Background())
	if err != nil {
		t.Fatalf("Voices().GetDefaultSettings() error = %v", err)
	}
	if settings == nil {
		t.Fatal("Voices().GetDefaultSettings() returned nil")
	}
}

func TestVoicesGetValidation(t *testing.T) {
	client, _ := NewClient()

	_, err := client.Voices().Get(context.Background(), "")
	if err != ErrEmptyVoiceID {
		t.Errorf("Get('') error = %v, want %v", err, ErrEmptyVoiceID)
	}

	_, err = client.Voices().GetSettings(context.Background(), "")
	if err != ErrEmptyVoiceID {
		t.Errorf("GetSettings('') error = %v, want %v", err, ErrEmptyVoiceID)
	}
}
