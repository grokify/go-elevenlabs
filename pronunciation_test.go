package elevenlabs

import (
	"context"
	"os"
	"testing"
)

func TestPronunciationService(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping live test")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()

	// Test that service is accessible
	if client.Pronunciation() == nil {
		t.Error("Pronunciation() returned nil")
	}

	t.Run("List", func(t *testing.T) {
		resp, err := client.Pronunciation().List(ctx, nil)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}
		if resp == nil {
			t.Error("List() returned nil")
			return
		}
		// Dictionaries may be empty, that's OK
		t.Logf("Found %d pronunciation dictionaries", len(resp.Dictionaries))
	})

	t.Run("Get with empty ID", func(t *testing.T) {
		_, err := client.Pronunciation().Get(ctx, "")
		if err == nil {
			t.Error("Get() with empty ID should return error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("Get() with empty ID should return ValidationError, got %T", err)
		}
	})

	t.Run("Create with empty name", func(t *testing.T) {
		_, err := client.Pronunciation().Create(ctx, &CreatePronunciationDictionaryRequest{
			Name: "",
		})
		if err == nil {
			t.Error("Create() with empty name should return error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("Create() with empty name should return ValidationError, got %T", err)
		}
	})

	t.Run("RemoveRules with empty ID", func(t *testing.T) {
		err := client.Pronunciation().RemoveRules(ctx, "", []string{"test"})
		if err == nil {
			t.Error("RemoveRules() with empty ID should return error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("RemoveRules() with empty ID should return ValidationError, got %T", err)
		}
	})

	t.Run("RemoveRules with empty rules", func(t *testing.T) {
		err := client.Pronunciation().RemoveRules(ctx, "test-id", []string{})
		if err == nil {
			t.Error("RemoveRules() with empty rules should return error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("RemoveRules() with empty rules should return ValidationError, got %T", err)
		}
	})

	t.Run("Rename with empty ID", func(t *testing.T) {
		err := client.Pronunciation().Rename(ctx, "", "new-name")
		if err == nil {
			t.Error("Rename() with empty ID should return error")
		}
	})

	t.Run("Rename with empty name", func(t *testing.T) {
		err := client.Pronunciation().Rename(ctx, "test-id", "")
		if err == nil {
			t.Error("Rename() with empty name should return error")
		}
	})

	t.Run("Archive with empty ID", func(t *testing.T) {
		err := client.Pronunciation().Archive(ctx, "")
		if err == nil {
			t.Error("Archive() with empty ID should return error")
		}
	})
}
