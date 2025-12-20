package elevenlabs

import (
	"context"
	"testing"
)

func TestModelsList_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	models, err := client.Models().List(context.Background())
	if err != nil {
		t.Fatalf("Models().List() error = %v", err)
	}
	if len(models) == 0 {
		t.Error("Models().List() returned empty list")
	}

	// Check that models have required fields
	for _, m := range models {
		if m.ModelID == "" {
			t.Error("Model has empty ModelID")
		}
		if m.Name == "" {
			t.Error("Model has empty Name")
		}
	}
}

func TestModelsListTTSModels_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	models, err := client.Models().ListTTSModels(context.Background())
	if err != nil {
		t.Fatalf("Models().ListTTSModels() error = %v", err)
	}
	if len(models) == 0 {
		t.Error("Models().ListTTSModels() returned empty list")
	}

	// Verify all returned models support TTS
	for _, m := range models {
		if !m.CanDoTextToSpeech {
			t.Errorf("Model %s does not support TTS but was returned by ListTTSModels", m.ModelID)
		}
	}
}

func TestModelsContainsDefaultModel_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	models, err := client.Models().List(context.Background())
	if err != nil {
		t.Fatalf("Models().List() error = %v", err)
	}

	found := false
	for _, m := range models {
		if m.ModelID == DefaultModelID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Default model %s not found in models list", DefaultModelID)
	}
}
