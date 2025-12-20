package voices

import (
	"testing"
)

func TestPremadeVoices(t *testing.T) {
	voices := PremadeVoices()
	if len(voices) == 0 {
		t.Error("PremadeVoices should return voices")
	}

	// Check that all voices have required fields
	for _, v := range voices {
		if v.ID == "" {
			t.Errorf("Voice %s has empty ID", v.Name)
		}
		if v.Name == "" {
			t.Error("Voice has empty name")
		}
		if v.Gender == "" {
			t.Errorf("Voice %s has empty gender", v.Name)
		}
	}
}

func TestGetVoice(t *testing.T) {
	// Test finding Rachel
	v := GetVoice(Rachel)
	if v == nil {
		t.Fatal("GetVoice should find Rachel")
	}
	if v.Name != "Rachel" {
		t.Errorf("expected name 'Rachel', got '%s'", v.Name)
	}

	// Test not found
	v = GetVoice("nonexistent")
	if v != nil {
		t.Error("GetVoice should return nil for nonexistent ID")
	}
}

func TestGetVoiceByName(t *testing.T) {
	// Test case-insensitive lookup
	v := GetVoiceByName("rachel")
	if v == nil {
		t.Fatal("GetVoiceByName should find rachel (lowercase)")
	}
	if v.ID != Rachel {
		t.Errorf("expected ID %s, got %s", Rachel, v.ID)
	}

	v = GetVoiceByName("RACHEL")
	if v == nil {
		t.Fatal("GetVoiceByName should find RACHEL (uppercase)")
	}

	// Test not found
	v = GetVoiceByName("nonexistent")
	if v != nil {
		t.Error("GetVoiceByName should return nil for nonexistent name")
	}
}

func TestFilterByGender(t *testing.T) {
	females := FilterByGender("female")
	if len(females) == 0 {
		t.Error("FilterByGender should find female voices")
	}
	for _, v := range females {
		if v.Gender != "female" {
			t.Errorf("expected female, got %s for %s", v.Gender, v.Name)
		}
	}

	males := FilterByGender("male")
	if len(males) == 0 {
		t.Error("FilterByGender should find male voices")
	}

	nonBinary := FilterByGender("non-binary")
	if len(nonBinary) == 0 {
		t.Error("FilterByGender should find non-binary voices")
	}
}

func TestFilterByAccent(t *testing.T) {
	british := FilterByAccent("British")
	if len(british) == 0 {
		t.Error("FilterByAccent should find British voices")
	}
	for _, v := range british {
		if !containsFold(v.Accent, "British") {
			t.Errorf("expected British accent, got %s for %s", v.Accent, v.Name)
		}
	}

	american := FilterByAccent("American")
	if len(american) == 0 {
		t.Error("FilterByAccent should find American voices")
	}
}

func TestFilterByAge(t *testing.T) {
	young := FilterByAge("young")
	if len(young) == 0 {
		t.Error("FilterByAge should find young voices")
	}

	middleAged := FilterByAge("middle-aged")
	if len(middleAged) == 0 {
		t.Error("FilterByAge should find middle-aged voices")
	}

	old := FilterByAge("old")
	if len(old) == 0 {
		t.Error("FilterByAge should find old voices")
	}
}

func TestVoiceConstants(t *testing.T) {
	// Verify some key constants match the expected IDs
	tests := []struct {
		constant string
		expected string
	}{
		{Rachel, "21m00Tcm4TlvDq8ikWAM"},
		{Adam, "pNInz6obpgDQGcFmaJgB"},
		{Antoni, "ErXwobaYiN019PkySvjV"},
		{Josh, "TxGEqnHWrfWFTfGW9XjX"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("constant mismatch: got %s, expected %s", tt.constant, tt.expected)
		}
	}
}
