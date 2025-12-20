package elevenlabs

import (
	"context"
	"strings"
	"testing"
)

func TestVoiceDesignRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Generate sample text of 100+ characters
	sampleText := strings.Repeat("This is a sample text for voice preview. ", 5)

	tests := []struct {
		name    string
		req     *VoiceDesignRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty gender",
			req:     &VoiceDesignRequest{Age: VoiceAgeYoung, Accent: VoiceAccentAmerican, Text: sampleText},
			wantErr: true,
			errMsg:  "gender",
		},
		{
			name:    "empty age",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Accent: VoiceAccentAmerican, Text: sampleText},
			wantErr: true,
			errMsg:  "age",
		},
		{
			name:    "empty accent",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Age: VoiceAgeYoung, Text: sampleText},
			wantErr: true,
			errMsg:  "accent",
		},
		{
			name:    "empty text",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Age: VoiceAgeYoung, Accent: VoiceAccentAmerican},
			wantErr: true,
			errMsg:  "text",
		},
		{
			name:    "text too short",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Age: VoiceAgeYoung, Accent: VoiceAccentAmerican, Text: "Too short"},
			wantErr: true,
			errMsg:  "text",
		},
		{
			name:    "invalid accent strength too low",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Age: VoiceAgeYoung, Accent: VoiceAccentAmerican, Text: sampleText, AccentStrength: 0.1},
			wantErr: true,
			errMsg:  "accent_strength",
		},
		{
			name:    "invalid accent strength too high",
			req:     &VoiceDesignRequest{Gender: VoiceGenderFemale, Age: VoiceAgeYoung, Accent: VoiceAccentAmerican, Text: sampleText, AccentStrength: 3.0},
			wantErr: true,
			errMsg:  "accent_strength",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.VoiceDesign().GeneratePreview(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("GeneratePreview() should return error")
					return
				}
				var valErr *ValidationError
				if isValidationError(err, &valErr) {
					if !strings.Contains(valErr.Field, tt.errMsg) {
						t.Errorf("ValidationError field = %s, want to contain %s", valErr.Field, tt.errMsg)
					}
				}
			}
		})
	}
}

func TestVoiceDesignService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.VoiceDesign() == nil {
		t.Error("VoiceDesign() returned nil")
	}
}

func TestVoiceGenderConstants(t *testing.T) {
	if VoiceGenderFemale != "female" {
		t.Errorf("VoiceGenderFemale = %s, want female", VoiceGenderFemale)
	}
	if VoiceGenderMale != "male" {
		t.Errorf("VoiceGenderMale = %s, want male", VoiceGenderMale)
	}
}

func TestVoiceAgeConstants(t *testing.T) {
	if VoiceAgeYoung != "young" {
		t.Errorf("VoiceAgeYoung = %s, want young", VoiceAgeYoung)
	}
	if VoiceAgeMiddleAged != "middle_aged" {
		t.Errorf("VoiceAgeMiddleAged = %s, want middle_aged", VoiceAgeMiddleAged)
	}
	if VoiceAgeOld != "old" {
		t.Errorf("VoiceAgeOld = %s, want old", VoiceAgeOld)
	}
}

func TestVoiceAccentConstants(t *testing.T) {
	if VoiceAccentBritish != "british" {
		t.Errorf("VoiceAccentBritish = %s, want british", VoiceAccentBritish)
	}
	if VoiceAccentAmerican != "american" {
		t.Errorf("VoiceAccentAmerican = %s, want american", VoiceAccentAmerican)
	}
	if VoiceAccentAfrican != "african" {
		t.Errorf("VoiceAccentAfrican = %s, want african", VoiceAccentAfrican)
	}
	if VoiceAccentAustralian != "australian" {
		t.Errorf("VoiceAccentAustralian = %s, want australian", VoiceAccentAustralian)
	}
	if VoiceAccentIndian != "indian" {
		t.Errorf("VoiceAccentIndian = %s, want indian", VoiceAccentIndian)
	}
}

func TestSaveVoiceRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test empty generated_voice_id
	_, err := client.VoiceDesign().SaveVoice(ctx, &SaveVoiceRequest{VoiceName: "Test"})
	if err == nil {
		t.Error("SaveVoice() with empty generated_voice_id should return error")
	}

	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "generated_voice_id" {
		t.Errorf("ValidationError field = %s, want generated_voice_id", valErr.Field)
	}

	// Test empty voice_name
	_, err = client.VoiceDesign().SaveVoice(ctx, &SaveVoiceRequest{GeneratedVoiceID: "test-id"})
	if err == nil {
		t.Error("SaveVoice() with empty voice_name should return error")
	}

	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "voice_name" {
		t.Errorf("ValidationError field = %s, want voice_name", valErr.Field)
	}
}

func TestSimpleVoiceDesign(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test with text too short (validation will fail)
	_, err := client.VoiceDesign().Simple(ctx, VoiceGenderFemale, VoiceAgeYoung, VoiceAccentAmerican, "short")
	if err == nil {
		t.Error("Simple() with short text should return error")
	}
}
