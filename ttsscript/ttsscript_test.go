package ttsscript

import (
	"strings"
	"testing"
)

func TestParseScript(t *testing.T) {
	jsonData := `{
		"title": "Test Script",
		"default_voices": {"en": "voice-1"},
		"pronunciations": {"API": {"en": "A P I"}},
		"slides": [
			{
				"title": "Intro",
				"segments": [
					{"text": {"en": "Hello API world"}, "pause_after": "500ms"}
				]
			}
		]
	}`

	script, err := ParseScript([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseScript failed: %v", err)
	}

	if script.Title != "Test Script" {
		t.Errorf("expected title 'Test Script', got '%s'", script.Title)
	}

	if len(script.Slides) != 1 {
		t.Errorf("expected 1 slide, got %d", len(script.Slides))
	}

	if script.Slides[0].Title != "Intro" {
		t.Errorf("expected slide title 'Intro', got '%s'", script.Slides[0].Title)
	}
}

func TestCompiler(t *testing.T) {
	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Pronunciations: map[string]map[string]string{
			"API": {"en": "A P I"},
		},
		Slides: []Slide{
			{
				Title: "Slide 1",
				Segments: []Segment{
					{
						Text:       map[string]string{"en": "Hello API world"},
						PauseAfter: "500ms",
					},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	seg := segments[0]

	// Check pronunciation was applied
	if seg.Text != "Hello A P I world" {
		t.Errorf("expected 'Hello A P I world', got '%s'", seg.Text)
	}

	// Check voice was set
	if seg.VoiceID != "voice-1" {
		t.Errorf("expected voice 'voice-1', got '%s'", seg.VoiceID)
	}

	// Check pause was parsed (500ms + default slide pause 800ms = 800ms since slide pause > segment pause)
	if seg.PauseAfterMs != 800 {
		t.Errorf("expected pause 800ms, got %dms", seg.PauseAfterMs)
	}
}

func TestSSMLFormatter(t *testing.T) {
	segments := []CompiledSegment{
		{
			SlideIndex:   0,
			SegmentIndex: 0,
			SlideTitle:   "Intro",
			Text:         "Hello world",
			PauseAfterMs: 500,
		},
	}

	formatter := NewSSMLFormatter()
	ssml := formatter.Format(segments, "en")

	if !strings.Contains(ssml, "<speak") {
		t.Error("SSML should contain <speak> tag")
	}

	if !strings.Contains(ssml, "Hello world") {
		t.Error("SSML should contain the text")
	}

	if !strings.Contains(ssml, `<break time="500ms"`) {
		t.Error("SSML should contain break element")
	}

	if !strings.Contains(ssml, "<!-- Slide 1: Intro -->") {
		t.Error("SSML should contain slide comment")
	}
}

func TestElevenLabsFormatter(t *testing.T) {
	segments := []CompiledSegment{
		{
			SlideIndex:   0,
			SegmentIndex: 0,
			SlideTitle:   "Intro",
			Text:         "Hello world",
			VoiceID:      "voice-1",
			PauseAfterMs: 500,
		},
		{
			SlideIndex:   0,
			SegmentIndex: 1,
			Text:         "Second segment",
			VoiceID:      "voice-2",
		},
	}

	formatter := NewElevenLabsFormatter()
	result := formatter.Format(segments)

	if len(result) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(result))
	}

	if result[0].VoiceID != "voice-1" {
		t.Errorf("expected voice 'voice-1', got '%s'", result[0].VoiceID)
	}

	if result[0].PauseAfterMs != 500 {
		t.Errorf("expected pause 500ms, got %dms", result[0].PauseAfterMs)
	}

	// Test grouping by voice
	groups := formatter.GroupByVoice(result)
	if len(groups) != 2 {
		t.Errorf("expected 2 voice groups, got %d", len(groups))
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"500ms", 500},
		{"1s", 1000},
		{"1.5s", 1500},
		{"2s", 2000},
		{"", 0},
		{"100ms", 100},
	}

	for _, tt := range tests {
		result := ParseDuration(tt.input)
		if result != tt.expected {
			t.Errorf("ParseDuration(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{500, "500ms"},
		{1000, "1s"},
		{2000, "2s"},
		{1500, "1500ms"},
		{0, ""},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDuration(%d) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestEscapeSSML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & world", "Hello &amp; world"},
		{"<tag>", "&lt;tag&gt;"},
		{`Say "hello"`, `Say &quot;hello&quot;`},
		{"It's fine", "It&apos;s fine"},
	}

	for _, tt := range tests {
		result := EscapeSSML(tt.input)
		if result != tt.expected {
			t.Errorf("EscapeSSML(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestScriptLanguages(t *testing.T) {
	script := &Script{
		Slides: []Slide{
			{
				Segments: []Segment{
					{Text: map[string]string{"en": "Hello", "es": "Hola"}},
					{Text: map[string]string{"en": "World", "fr": "Monde"}},
				},
			},
		},
	}

	langs := script.Languages()
	if len(langs) != 3 {
		t.Errorf("expected 3 languages, got %d", len(langs))
	}

	langSet := make(map[string]bool)
	for _, l := range langs {
		langSet[l] = true
	}

	for _, expected := range []string{"en", "es", "fr"} {
		if !langSet[expected] {
			t.Errorf("expected language %q in result", expected)
		}
	}
}

func TestScriptValidate(t *testing.T) {
	// Valid script
	valid := &Script{
		Slides: []Slide{
			{Segments: []Segment{{Text: map[string]string{"en": "Hello"}}}},
		},
	}
	if issues := valid.Validate(); len(issues) != 0 {
		t.Errorf("valid script should have no issues, got: %v", issues)
	}

	// Empty script
	empty := &Script{}
	if issues := empty.Validate(); len(issues) == 0 {
		t.Error("empty script should have issues")
	}

	// Slide with no segments
	noSegs := &Script{
		Slides: []Slide{{}},
	}
	if issues := noSegs.Validate(); len(issues) == 0 {
		t.Error("slide with no segments should have issues")
	}
}
