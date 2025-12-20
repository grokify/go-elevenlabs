package elevenlabs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPronunciationRuleValidate(t *testing.T) {
	tests := []struct {
		name    string
		rule    PronunciationRule
		wantErr bool
	}{
		{
			name:    "empty grapheme",
			rule:    PronunciationRule{Grapheme: "", Alias: "test"},
			wantErr: true,
		},
		{
			name:    "no alias or phoneme",
			rule:    PronunciationRule{Grapheme: "ADK"},
			wantErr: true,
		},
		{
			name:    "both alias and phoneme",
			rule:    PronunciationRule{Grapheme: "ADK", Alias: "test", Phoneme: "test"},
			wantErr: true,
		},
		{
			name:    "valid with alias",
			rule:    PronunciationRule{Grapheme: "ADK", Alias: "Agent Development Kit"},
			wantErr: false,
		},
		{
			name:    "valid with phoneme",
			rule:    PronunciationRule{Grapheme: "nginx", Phoneme: "ˈɛndʒɪnˈɛks"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRulesFromMap(t *testing.T) {
	m := map[string]string{
		"ADK":     "Agent Development Kit",
		"kubectl": "kube control",
	}

	rules := RulesFromMap(m)

	if len(rules) != 2 {
		t.Errorf("RulesFromMap() returned %d rules, want 2", len(rules))
	}

	// Check that all rules are alias-based
	for _, rule := range rules {
		if rule.Alias == "" {
			t.Errorf("RulesFromMap() rule for %s has no alias", rule.Grapheme)
		}
		if rule.Phoneme != "" {
			t.Errorf("RulesFromMap() rule for %s should not have phoneme", rule.Grapheme)
		}
	}
}

func TestPronunciationRulesToPLS(t *testing.T) {
	rules := PronunciationRules{
		{Grapheme: "ADK", Alias: "Agent Development Kit"},
		{Grapheme: "nginx", Phoneme: "ˈɛndʒɪnˈɛks"},
	}

	pls, err := rules.ToPLSString("en-US")
	if err != nil {
		t.Fatalf("ToPLSString() error = %v", err)
	}

	// Check XML structure
	if !strings.Contains(pls, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("ToPLSString() missing XML declaration")
	}
	if !strings.Contains(pls, `<lexicon`) {
		t.Error("ToPLSString() missing lexicon element")
	}
	if !strings.Contains(pls, `xmlns="http://www.w3.org/2005/01/pronunciation-lexicon"`) {
		t.Error("ToPLSString() missing xmlns attribute")
	}
	if !strings.Contains(pls, `xml:lang="en-US"`) {
		t.Error("ToPLSString() missing xml:lang attribute")
	}
	if !strings.Contains(pls, `<grapheme>ADK</grapheme>`) {
		t.Error("ToPLSString() missing ADK grapheme")
	}
	if !strings.Contains(pls, `<alias>Agent Development Kit</alias>`) {
		t.Error("ToPLSString() missing ADK alias")
	}
	if !strings.Contains(pls, `<grapheme>nginx</grapheme>`) {
		t.Error("ToPLSString() missing nginx grapheme")
	}
	if !strings.Contains(pls, `<phoneme>ˈɛndʒɪnˈɛks</phoneme>`) {
		t.Error("ToPLSString() missing nginx phoneme")
	}
}

func TestParseRulesFromJSON(t *testing.T) {
	jsonData := `[
		{"grapheme": "ADK", "alias": "Agent Development Kit"},
		{"grapheme": "kubectl", "alias": "kube control"}
	]`

	rules, err := ParseRulesFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseRulesFromJSON() error = %v", err)
	}

	if len(rules) != 2 {
		t.Errorf("ParseRulesFromJSON() returned %d rules, want 2", len(rules))
	}
}

func TestLoadRulesFromJSON(t *testing.T) {
	// Create a temp file with JSON content
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "rules.json")

	rules := []PronunciationRule{
		{Grapheme: "ADK", Alias: "Agent Development Kit"},
		{Grapheme: "API", Alias: "A P I"},
	}

	data, err := json.Marshal(rules)
	if err != nil {
		t.Fatalf("Failed to marshal rules: %v", err)
	}

	if err := os.WriteFile(jsonFile, data, 0600); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Load and verify
	loaded, err := LoadRulesFromJSON(jsonFile)
	if err != nil {
		t.Fatalf("LoadRulesFromJSON() error = %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("LoadRulesFromJSON() returned %d rules, want 2", len(loaded))
	}
}

func TestPronunciationRulesGraphemes(t *testing.T) {
	rules := PronunciationRules{
		{Grapheme: "ADK", Alias: "Agent Development Kit"},
		{Grapheme: "kubectl", Alias: "kube control"},
	}

	graphemes := rules.Graphemes()

	if len(graphemes) != 2 {
		t.Errorf("Graphemes() returned %d items, want 2", len(graphemes))
	}

	// Check values (order may vary)
	found := make(map[string]bool)
	for _, g := range graphemes {
		found[g] = true
	}
	if !found["ADK"] || !found["kubectl"] {
		t.Errorf("Graphemes() = %v, want [ADK, kubectl]", graphemes)
	}
}

func TestPronunciationRulesString(t *testing.T) {
	rules := PronunciationRules{
		{Grapheme: "ADK", Alias: "Agent Development Kit"},
		{Grapheme: "nginx", Phoneme: "ˈɛndʒɪnˈɛks"},
	}

	s := rules.String()

	if !strings.Contains(s, "ADK → Agent Development Kit") {
		t.Error("String() missing alias format")
	}
	if !strings.Contains(s, "nginx → [ˈɛndʒɪnˈɛks]") {
		t.Error("String() missing phoneme format")
	}
}

func TestPronunciationRulesSavePLS(t *testing.T) {
	rules := PronunciationRules{
		{Grapheme: "ADK", Alias: "Agent Development Kit"},
	}

	tmpDir := t.TempDir()
	plsFile := filepath.Join(tmpDir, "rules.pls")

	err := rules.SavePLS(plsFile, "en-US")
	if err != nil {
		t.Fatalf("SavePLS() error = %v", err)
	}

	// Verify file was created and contains expected content
	data, err := os.ReadFile(plsFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "<grapheme>ADK</grapheme>") {
		t.Error("SavePLS() file missing expected content")
	}
}

func TestParseRulesFromJSONInvalid(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "invalid JSON",
			json: `not valid json`,
		},
		{
			name: "missing grapheme",
			json: `[{"alias": "test"}]`,
		},
		{
			name: "missing alias and phoneme",
			json: `[{"grapheme": "ADK"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRulesFromJSON([]byte(tt.json))
			if err == nil {
				t.Error("ParseRulesFromJSON() expected error, got nil")
			}
		})
	}
}
