// Package ttsscript provides a structured format for authoring multilingual
// TTS (Text-to-Speech) scripts that can be compiled to various output formats.
//
// This package is engine-agnostic and can be used with any TTS provider including
// ElevenLabs, Google Cloud TTS, Amazon Polly, Azure TTS, and others.
//
// # Why Use ttsscript?
//
// Instead of storing raw SSML (which is engine-specific and hard to edit), store
// your scripts in a structured JSON format that:
//
//   - Supports multiple languages in a single file
//   - Handles pronunciations/acronyms separately from content
//   - Can be compiled to any TTS engine format
//   - Is easy to edit and version control
//
// # Basic Usage
//
// Create a script JSON file:
//
//	{
//	  "title": "My Course",
//	  "default_voices": {"en": "voice-id", "es": "voice-id-2"},
//	  "pronunciations": {
//	    "API": {"en": "A P I", "es": "A P I"},
//	    "SDK": {"en": "S D K", "es": "S D K"}
//	  },
//	  "slides": [
//	    {
//	      "title": "Introduction",
//	      "segments": [
//	        {
//	          "text": {"en": "Welcome to the API course", "es": "Bienvenidos al curso de API"},
//	          "pause_after": "500ms"
//	        }
//	      ]
//	    }
//	  ]
//	}
//
// Load and compile for ElevenLabs:
//
//	script, _ := ttsscript.LoadScript("script.json")
//	compiler := ttsscript.NewCompiler()
//	segments, _ := compiler.Compile(script, "en")
//
//	formatter := ttsscript.NewElevenLabsFormatter()
//	jobs := formatter.Format(segments)
//
//	for _, job := range jobs {
//	    // Generate TTS for each segment
//	    audio, _ := client.TextToSpeech().Simple(ctx, job.VoiceID, job.Text)
//	    // Save with pause information for post-processing
//	}
//
// Compile to SSML for Google TTS:
//
//	formatter := ttsscript.NewSSMLFormatter()
//	ssml, _ := formatter.FormatScript(script, "en")
//	// Use ssml with Google Cloud TTS API
//
// # Script Structure
//
// A Script contains:
//   - Metadata (title, description, default language)
//   - Default voices per language
//   - Global pronunciations
//   - Slides/sections containing segments
//
// Each Segment contains:
//   - Text in multiple languages
//   - Voice overrides per language
//   - Pause before/after
//   - Prosody settings (rate, pitch, emphasis)
//   - Segment-specific pronunciations
//
// # Compilation Process
//
// 1. Load the script from JSON
// 2. Create a Compiler and optionally add additional pronunciations
// 3. Compile for a specific language to get CompiledSegments
// 4. Format the segments for your target TTS engine
//
// # Formatters
//
// SSMLFormatter: Outputs W3C SSML compatible with Google, Amazon, Azure
// ElevenLabsFormatter: Outputs segments ready for ElevenLabs TTS API
//
// # Pronunciation Handling
//
// Pronunciations are applied at compile time with this priority:
// 1. Compiler-level (added via AddPronunciation)
// 2. Segment-level (in segment.pronunciations)
// 3. Script-level (in script.pronunciations)
//
// This allows overrides at any level. Terms are matched case-insensitively
// with word boundaries.
package ttsscript
