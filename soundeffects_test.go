package elevenlabs

import (
	"context"
	"os"
	"testing"
)

func TestSoundEffectRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *SoundEffectRequest
		wantErr bool
	}{
		{
			name:    "empty text",
			req:     &SoundEffectRequest{Text: ""},
			wantErr: true,
		},
		{
			name:    "valid request",
			req:     &SoundEffectRequest{Text: "car engine starting"},
			wantErr: false,
		},
		{
			name: "duration too short",
			req: &SoundEffectRequest{
				Text:            "thunder",
				DurationSeconds: 0.1,
			},
			wantErr: true,
		},
		{
			name: "duration too long",
			req: &SoundEffectRequest{
				Text:            "thunder",
				DurationSeconds: 35,
			},
			wantErr: true,
		},
		{
			name: "valid duration",
			req: &SoundEffectRequest{
				Text:            "thunder",
				DurationSeconds: 5,
			},
			wantErr: false,
		},
		{
			name: "prompt influence out of range",
			req: &SoundEffectRequest{
				Text:            "thunder",
				PromptInfluence: 1.5,
			},
			wantErr: true,
		},
		{
			name: "valid prompt influence",
			req: &SoundEffectRequest{
				Text:            "thunder",
				PromptInfluence: 0.5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSoundEffectsService(t *testing.T) {
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
	if client.SoundEffects() == nil {
		t.Error("SoundEffects() returned nil")
	}

	// Test Simple with short sound effect
	t.Run("Simple", func(t *testing.T) {
		audio, err := client.SoundEffects().Simple(ctx, "short beep")
		if err != nil {
			t.Errorf("Simple() error = %v", err)
			return
		}
		if audio == nil {
			t.Error("Simple() returned nil audio")
		}
	})
}
