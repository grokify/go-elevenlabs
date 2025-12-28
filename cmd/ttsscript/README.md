# ttsscript

A CLI tool for generating TTS audio from JSON script files using ElevenLabs.

## Overview

`ttsscript` reads a structured JSON script file and generates audio files using the ElevenLabs TTS API. It supports:

- Multilingual scripts with per-language voice configuration
- Section headers with spoken titles
- Custom pronunciation rules
- Per-segment or per-slide audio output
- Manifest generation for video editing workflows

## Installation

```bash
go install github.com/agentplexus/go-elevenlabs/cmd/ttsscript@latest
```

Or build from source:

```bash
git clone https://github.com/agentplexus/go-elevenlabs.git
cd go-elevenlabs
go build -o ttsscript ./cmd/ttsscript
```

## Requirements

- **ElevenLabs API key**: Set via `ELEVENLABS_API_KEY` environment variable
- **ffmpeg** (optional): Required only for `--per-slide` mode

## Usage

```bash
ttsscript [flags] <script.json>
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-lang` | `en` | Language code to generate (must exist in script) |
| `-output` | `./output` | Output directory for audio files |
| `-per-slide` | `false` | Concatenate segments into per-slide audio files |
| `-manifest` | `true` | Generate manifest JSON file |
| `-dry-run` | `false` | Preview output without calling API |
| `-model` | `eleven_multilingual_v2` | ElevenLabs model ID |

### Examples

```bash
# Preview what would be generated
ttsscript -dry-run script.json

# Generate English audio
ttsscript -lang en -output ./audio script.json

# Generate Spanish audio with per-slide output
ttsscript -lang es -output ./audio -per-slide script.json

# Use a specific model
ttsscript -model eleven_turbo_v2_5 script.json
```

## Script Format

Scripts are JSON files with the following structure:

```json
{
  "title": "Course Introduction",
  "description": "An introduction to the course",
  "default_language": "en",
  "default_voices": {
    "en": "21m00Tcm4TlvDq8ikWAM",
    "es": "EXAVITQu4vr4xnSDxMaL"
  },
  "pronunciations": {
    "API": {"en": "A P I", "es": "A P I"},
    "SDK": {"en": "S D K", "es": "S D K"}
  },
  "slides": [
    {
      "title": "Welcome",
      "is_section_header": true,
      "segments": [
        {
          "text": {
            "en": "Welcome to this course.",
            "es": "Bienvenidos a este curso."
          },
          "pause_after": "500ms"
        }
      ]
    }
  ]
}
```

### Script Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Script title (metadata) |
| `description` | string | Script description (metadata) |
| `default_language` | string | Primary language code |
| `default_voices` | object | Map of language code to ElevenLabs voice ID |
| `pronunciations` | object | Global pronunciation rules (term → language → replacement) |
| `slides` | array | Ordered list of slides |

### Slide Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Slide title |
| `notes` | string | Speaker notes (not rendered to audio) |
| `is_section_header` | bool | Marks slide as section start |
| `speak_title` | bool | Speak title before segments (default: true for section headers) |
| `title_voice` | object | Voice override for title by language |
| `title_pause_after` | string | Pause after title (default: 500ms for sections, 300ms otherwise) |
| `segments` | array | Audio segments for this slide |

### Segment Fields

| Field | Type | Description |
|-------|------|-------------|
| `text` | object | Text by language code (required) |
| `voice` | object | Voice override by language |
| `pause_before` | string | Pause before segment (e.g., "500ms", "1s") |
| `pause_after` | string | Pause after segment |
| `emphasis` | string | Emphasis level: "strong", "moderate", "reduced" |
| `rate` | string | Speaking rate: "slow", "medium", "fast", or percentage |
| `pitch` | string | Pitch adjustment: "low", "medium", "high", or percentage |
| `pronunciations` | object | Segment-specific pronunciation overrides |

## Output Structure

### Per-Segment Mode (default)

```
output/
├── slide01_title_en.mp3      # Section header title
├── slide01_seg01_en.mp3      # First segment
├── slide01_seg02_en.mp3      # Second segment
├── slide02_seg01_en.mp3      # Next slide's segment
└── manifest_en.json          # Generation manifest
```

### Per-Slide Mode (`--per-slide`)

```
output/
├── slide01_title_en.mp3      # Individual segments (kept)
├── slide01_seg01_en.mp3
├── slide01_seg02_en.mp3
├── slide01_en.mp3            # Concatenated slide audio
├── slide02_seg01_en.mp3
├── slide02_en.mp3
└── manifest_en.json
```

## Manifest Format

The manifest file tracks all generated segments for downstream processing:

```json
[
  {
    "slide_index": 0,
    "segment_index": -1,
    "slide_title": "Introduction",
    "is_title_segment": true,
    "is_section_header": true,
    "text": "Introduction",
    "voice_id": "21m00Tcm4TlvDq8ikWAM",
    "language": "en",
    "output_file": "./output/slide01_title_en.mp3",
    "pause_before_ms": 0,
    "pause_after_ms": 500
  },
  {
    "slide_index": 0,
    "segment_index": 0,
    "slide_title": "Introduction",
    "text": "Welcome to the course.",
    "voice_id": "21m00Tcm4TlvDq8ikWAM",
    "language": "en",
    "output_file": "./output/slide01_seg01_en.mp3",
    "pause_after_ms": 800
  }
]
```

## Example Script

Here's a complete example script:

```json
{
  "title": "Go Programming Introduction",
  "default_voices": {
    "en": "21m00Tcm4TlvDq8ikWAM"
  },
  "pronunciations": {
    "Go": {"en": "Go"},
    "goroutine": {"en": "go routine"}
  },
  "slides": [
    {
      "title": "Introduction",
      "is_section_header": true,
      "segments": [
        {
          "text": {"en": "Welcome to this introduction to Go programming."},
          "pause_after": "800ms"
        },
        {
          "text": {"en": "Go is a fast, simple, and powerful language."},
          "pause_after": "500ms"
        }
      ]
    },
    {
      "title": "Key Features",
      "segments": [
        {
          "text": {"en": "Go compiles directly to machine code."},
          "pause_after": "300ms"
        },
        {
          "text": {"en": "It has excellent support for goroutine-based concurrency."},
          "pause_after": "300ms"
        }
      ]
    }
  ]
}
```

## LMS Video Workflow

For Learning Management System (LMS) video production:

1. **Generate per-segment audio** (default mode):
   ```bash
   ttsscript -lang en -output ./audio script.json
   ```

2. **Import manifest into video editor** - Use `manifest_en.json` for timing info

3. **Align segments to slides** - Match `slide_index` to your slide deck

4. **Or use per-slide mode** for simpler workflows:
   ```bash
   ttsscript -lang en -output ./audio -per-slide script.json
   ```

The per-segment approach gives you maximum flexibility for timing adjustments and re-recording individual segments without regenerating entire slides.

## Troubleshooting

### "ELEVENLABS_API_KEY environment variable is required"

Set your API key:
```bash
export ELEVENLABS_API_KEY=your_api_key_here
```

### "ffmpeg is required for --per-slide mode"

Install ffmpeg:
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows (with Chocolatey)
choco install ffmpeg
```

### "no voice ID configured"

Ensure your script has `default_voices` set for the language you're generating, or each segment has a `voice` override.

## See Also

- [ElevenLabs API Documentation](https://elevenlabs.io/docs)
- [ElevenLabs Voice Library](https://elevenlabs.io/voice-library)
- [ttsscript package documentation](../../ttsscript/)
