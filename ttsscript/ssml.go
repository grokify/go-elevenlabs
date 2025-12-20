package ttsscript

import (
	"fmt"
	"strings"
)

// SSMLFormatter formats compiled segments as SSML.
// Compatible with Google Cloud TTS, Amazon Polly, Azure TTS, and others.
type SSMLFormatter struct {
	// Version is the SSML version (default: "1.1").
	Version string

	// IncludeComments includes slide title comments in output.
	IncludeComments bool

	// IndentSpaces is the number of spaces for indentation.
	IndentSpaces int
}

// NewSSMLFormatter creates a new SSML formatter with default settings.
func NewSSMLFormatter() *SSMLFormatter {
	return &SSMLFormatter{
		Version:         "1.1",
		IncludeComments: true,
		IndentSpaces:    2,
	}
}

// Format formats compiled segments as SSML.
func (f *SSMLFormatter) Format(segments []CompiledSegment, language string) string {
	var sb strings.Builder
	indent := strings.Repeat(" ", f.IndentSpaces)

	// SSML header
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`<speak version="%s" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="%s">`,
		f.Version, language))
	sb.WriteString("\n")

	currentSlide := -1
	for _, seg := range segments {
		// Add slide comment when slide changes
		if f.IncludeComments && seg.SlideIndex != currentSlide {
			currentSlide = seg.SlideIndex
			if seg.SlideTitle != "" {
				sb.WriteString(fmt.Sprintf("%s<!-- Slide %d: %s -->\n", indent, seg.SlideIndex+1, seg.SlideTitle))
			} else {
				sb.WriteString(fmt.Sprintf("%s<!-- Slide %d -->\n", indent, seg.SlideIndex+1))
			}
		}

		// Add pause before
		if seg.PauseBeforeMs > 0 {
			sb.WriteString(fmt.Sprintf(`%s<break time="%s"/>`, indent, FormatDuration(seg.PauseBeforeMs)))
			sb.WriteString("\n")
		}

		// Build the segment with optional prosody/emphasis
		f.writeSegmentContent(&sb, seg, indent)

		// Add pause after
		if seg.PauseAfterMs > 0 {
			sb.WriteString(fmt.Sprintf(`%s<break time="%s"/>`, indent, FormatDuration(seg.PauseAfterMs)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("</speak>\n")

	return sb.String()
}

// writeSegmentContent writes the segment content with prosody/emphasis wrappers.
func (f *SSMLFormatter) writeSegmentContent(sb *strings.Builder, seg CompiledSegment, indent string) {
	hasProsody := seg.Rate != "" || seg.Pitch != ""
	hasEmphasis := seg.Emphasis != ""

	sb.WriteString(indent)

	// Open prosody tag
	if hasProsody {
		sb.WriteString("<prosody")
		if seg.Rate != "" {
			sb.WriteString(fmt.Sprintf(` rate="%s"`, seg.Rate))
		}
		if seg.Pitch != "" {
			sb.WriteString(fmt.Sprintf(` pitch="%s"`, seg.Pitch))
		}
		sb.WriteString(">")
	}

	// Open emphasis tag
	if hasEmphasis {
		sb.WriteString(fmt.Sprintf(`<emphasis level="%s">`, seg.Emphasis))
	}

	// Write text content
	sb.WriteString(EscapeSSML(seg.Text))

	// Close emphasis tag
	if hasEmphasis {
		sb.WriteString("</emphasis>")
	}

	// Close prosody tag
	if hasProsody {
		sb.WriteString("</prosody>")
	}

	sb.WriteString("\n")
}

// FormatScript compiles and formats a script as SSML.
func (f *SSMLFormatter) FormatScript(script *Script, language string) (string, error) {
	compiler := NewCompiler()
	segments, err := compiler.Compile(script, language)
	if err != nil {
		return "", err
	}
	return f.Format(segments, language), nil
}

// EscapeSSML escapes special characters for SSML.
func EscapeSSML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// SSMLBreak generates an SSML break element.
func SSMLBreak(duration string) string {
	return fmt.Sprintf(`<break time="%s"/>`, duration)
}

// SSMLProsody wraps text in prosody tags.
func SSMLProsody(text, rate, pitch, volume string) string {
	var attrs []string
	if rate != "" {
		attrs = append(attrs, fmt.Sprintf(`rate="%s"`, rate))
	}
	if pitch != "" {
		attrs = append(attrs, fmt.Sprintf(`pitch="%s"`, pitch))
	}
	if volume != "" {
		attrs = append(attrs, fmt.Sprintf(`volume="%s"`, volume))
	}
	if len(attrs) == 0 {
		return text
	}
	return fmt.Sprintf("<prosody %s>%s</prosody>", strings.Join(attrs, " "), text)
}

// SSMLEmphasis wraps text in emphasis tags.
func SSMLEmphasis(text, level string) string {
	return fmt.Sprintf(`<emphasis level="%s">%s</emphasis>`, level, text)
}

// SSMLSayAs wraps text in say-as tags for specific interpretation.
func SSMLSayAs(text, interpretAs, format string) string {
	if format != "" {
		return fmt.Sprintf(`<say-as interpret-as="%s" format="%s">%s</say-as>`, interpretAs, format, text)
	}
	return fmt.Sprintf(`<say-as interpret-as="%s">%s</say-as>`, interpretAs, text)
}

// SSMLPhoneme wraps text with phonetic pronunciation.
func SSMLPhoneme(text, alphabet, ph string) string {
	return fmt.Sprintf(`<phoneme alphabet="%s" ph="%s">%s</phoneme>`, alphabet, ph, text)
}

// SSMLSub provides an alias for a word.
func SSMLSub(text, alias string) string {
	return fmt.Sprintf(`<sub alias="%s">%s</sub>`, alias, text)
}
