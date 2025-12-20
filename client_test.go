package elevenlabs

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test creating client without API key (uses environment)
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	// Verify services are initialized
	if client.TextToSpeech() == nil {
		t.Error("TextToSpeech() service is nil")
	}
	if client.Voices() == nil {
		t.Error("Voices() service is nil")
	}
	if client.Models() == nil {
		t.Error("Models() service is nil")
	}
	if client.History() == nil {
		t.Error("History() service is nil")
	}
	if client.User() == nil {
		t.Error("User() service is nil")
	}
	if client.Dubbing() == nil {
		t.Error("Dubbing() service is nil")
	}
	if client.SoundEffects() == nil {
		t.Error("SoundEffects() service is nil")
	}
	if client.Pronunciation() == nil {
		t.Error("Pronunciation() service is nil")
	}
	if client.Projects() == nil {
		t.Error("Projects() service is nil")
	}
	if client.SpeechToText() == nil {
		t.Error("SpeechToText() service is nil")
	}
	if client.ForcedAlignment() == nil {
		t.Error("ForcedAlignment() service is nil")
	}
	if client.AudioIsolation() == nil {
		t.Error("AudioIsolation() service is nil")
	}
	if client.TextToDialogue() == nil {
		t.Error("TextToDialogue() service is nil")
	}
	if client.VoiceDesign() == nil {
		t.Error("VoiceDesign() service is nil")
	}
	if client.Music() == nil {
		t.Error("Music() service is nil")
	}
	if client.API() == nil {
		t.Error("API() returned nil")
	}
}

func TestNewClientWithAPIKey(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("NewClient(WithAPIKey()) error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client, err := NewClient(
		WithAPIKey("test-api-key"),
		WithBaseURL("https://custom.api.com"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("baseURL = %s, want https://custom.api.com", client.baseURL)
	}
}

// Helper function to get API key for live tests
func getAPIKey(t *testing.T) string {
	t.Helper()
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping live API test")
	}
	return apiKey
}
