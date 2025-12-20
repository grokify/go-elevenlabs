package elevenlabs

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// PronunciationRule defines how a word or phrase should be pronounced.
// Rules can use either an alias (text substitution) or IPA phonemes.
//
// Example JSON:
//
//	[
//	  {"grapheme": "ADK", "alias": "Agent Development Kit"},
//	  {"grapheme": "kubectl", "alias": "kube control"},
//	  {"grapheme": "nginx", "phoneme": "ˈɛndʒɪnˈɛks"}
//	]
type PronunciationRule struct {
	// Grapheme is the text to match (required).
	Grapheme string `json:"grapheme"`

	// Alias is the replacement text (mutually exclusive with Phoneme).
	// This is the easier option - just specify what text should be read instead.
	Alias string `json:"alias,omitempty"`

	// Phoneme is the IPA pronunciation (mutually exclusive with Alias).
	// Use this for precise phonetic control.
	Phoneme string `json:"phoneme,omitempty"`
}

// Validate checks that the rule is valid.
func (r *PronunciationRule) Validate() error {
	if r.Grapheme == "" {
		return &ValidationError{Field: "grapheme", Message: "cannot be empty"}
	}
	if r.Alias == "" && r.Phoneme == "" {
		return &ValidationError{Field: "alias/phoneme", Message: "either alias or phoneme must be specified"}
	}
	if r.Alias != "" && r.Phoneme != "" {
		return &ValidationError{Field: "alias/phoneme", Message: "cannot specify both alias and phoneme"}
	}
	return nil
}

// PronunciationRules is a collection of pronunciation rules.
type PronunciationRules []PronunciationRule

// LoadRulesFromJSON loads pronunciation rules from a JSON file.
//
// Example file content:
//
//	[
//	  {"grapheme": "ADK", "alias": "Agent Development Kit"},
//	  {"grapheme": "API", "alias": "A P I"},
//	  {"grapheme": "SQL", "alias": "sequel"}
//	]
func LoadRulesFromJSON(filename string) (PronunciationRules, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading pronunciation rules file: %w", err)
	}

	var rules PronunciationRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parsing pronunciation rules JSON: %w", err)
	}

	// Validate all rules
	for i, rule := range rules {
		if err := rule.Validate(); err != nil {
			return nil, fmt.Errorf("rule %d: %w", i, err)
		}
	}

	return rules, nil
}

// ParseRulesFromJSON parses pronunciation rules from JSON bytes.
func ParseRulesFromJSON(data []byte) (PronunciationRules, error) {
	var rules PronunciationRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parsing pronunciation rules JSON: %w", err)
	}

	// Validate all rules
	for i, rule := range rules {
		if err := rule.Validate(); err != nil {
			return nil, fmt.Errorf("rule %d: %w", i, err)
		}
	}

	return rules, nil
}

// RulesFromMap creates pronunciation rules from a simple map.
// All entries are treated as alias substitutions.
//
// Example:
//
//	rules := RulesFromMap(map[string]string{
//	    "ADK":     "Agent Development Kit",
//	    "kubectl": "kube control",
//	    "API":     "A P I",
//	})
func RulesFromMap(m map[string]string) PronunciationRules {
	rules := make(PronunciationRules, 0, len(m))
	for grapheme, alias := range m {
		rules = append(rules, PronunciationRule{
			Grapheme: grapheme,
			Alias:    alias,
		})
	}
	return rules
}

// ToPLS converts the rules to PLS (Pronunciation Lexicon Specification) XML format.
// This is the format required by ElevenLabs API.
func (rules PronunciationRules) ToPLS(language string) ([]byte, error) {
	if language == "" {
		language = "en-US"
	}

	var lexemes []plsLexeme
	for _, rule := range rules {
		lexeme := plsLexeme{
			Grapheme: rule.Grapheme,
		}
		if rule.Alias != "" {
			lexeme.Alias = rule.Alias
		} else {
			lexeme.Phoneme = rule.Phoneme
		}
		lexemes = append(lexemes, lexeme)
	}

	lexicon := plsLexicon{
		Version:  "1.0",
		XMLNS:    "http://www.w3.org/2005/01/pronunciation-lexicon",
		Alphabet: "ipa",
		XMLLang:  language,
		Lexemes:  lexemes,
	}

	output, err := xml.MarshalIndent(lexicon, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("generating PLS XML: %w", err)
	}

	// Add XML declaration
	return []byte(xml.Header + string(output)), nil
}

// ToPLSString is a convenience method that returns the PLS as a string.
func (rules PronunciationRules) ToPLSString(language string) (string, error) {
	data, err := rules.ToPLS(language)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SavePLS writes the rules to a PLS file.
func (rules PronunciationRules) SavePLS(filename, language string) error {
	data, err := rules.ToPLS(language)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0600)
}

// Graphemes returns a list of all graphemes (the words being defined).
func (rules PronunciationRules) Graphemes() []string {
	result := make([]string, len(rules))
	for i, rule := range rules {
		result[i] = rule.Grapheme
	}
	return result
}

// String returns a human-readable summary of the rules.
func (rules PronunciationRules) String() string {
	var sb strings.Builder
	for _, rule := range rules {
		if rule.Alias != "" {
			sb.WriteString(fmt.Sprintf("%s → %s\n", rule.Grapheme, rule.Alias))
		} else {
			sb.WriteString(fmt.Sprintf("%s → [%s]\n", rule.Grapheme, rule.Phoneme))
		}
	}
	return sb.String()
}

// PLS XML structures (internal)

type plsLexicon struct {
	XMLName  xml.Name    `xml:"lexicon"`
	Version  string      `xml:"version,attr"`
	XMLNS    string      `xml:"xmlns,attr"`
	Alphabet string      `xml:"alphabet,attr"`
	XMLLang  string      `xml:"xml:lang,attr"`
	Lexemes  []plsLexeme `xml:"lexeme"`
}

type plsLexeme struct {
	Grapheme string `xml:"grapheme"`
	Alias    string `xml:"alias,omitempty"`
	Phoneme  string `xml:"phoneme,omitempty"`
}
