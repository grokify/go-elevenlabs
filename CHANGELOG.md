# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2024-12-28

### Added
- **WebSocket TTS** - Real-time text-to-speech streaming via WebSocket
  - Low-latency voice synthesis for LLM integration
  - Word-level alignment timestamps
  - Configurable latency optimization
- **WebSocket STT** - Real-time speech-to-text streaming
  - Partial/interim transcription results
  - Word-level timing with confidence scores
  - Automatic language detection
- **Speech-to-Speech** - Voice conversion service
  - Transform speech to different voices
  - Background noise removal
  - Streaming conversion support
- **Twilio Integration** - Phone call integration for voice agents
  - Incoming call registration
  - Outbound calls via Twilio and SIP
  - Dynamic variables for prompt injection
- **Phone Number Management** - Manage phone numbers for agents
- Added `gorilla/websocket` dependency for WebSocket support
- 4 new documentation pages for real-time services
- "Real-Time" section in MkDocs navigation
- Release notes documentation

### Changed
- Updated API coverage to ~75 methods
- Updated README with real-time service examples
- Updated presentation with new service counts

## [0.2.0] - 2024-12-28

### Changed
- **Repository Transfer** - Moved from `grokify/go-elevenlabs` to `agentplexus/go-elevenlabs`
- Updated all internal import paths
- Updated documentation URLs and references

### Added
- **MkDocs Documentation Site** - 28 pages with Material theme
- **`cmd/ttsscript`** - Command-line tool for TTS script processing
- LMS/Udemy course production guide
- Pronunciation rules guide
- TTS script authoring guide
- API coverage tracking with method-level details

## [0.1.0] - 2024-12-28

### Added
- Initial release of the ElevenLabs Go SDK
- **15 Service Wrappers**:
  - Text-to-Speech (streaming, timestamps)
  - Speech-to-Text (diarization)
  - Voices (management)
  - Voice Design
  - Sound Effects
  - Music (composition, stem separation)
  - Audio Isolation
  - Forced Alignment
  - Text-to-Dialogue
  - Projects (Studio)
  - Pronunciation Dictionaries
  - Dubbing
  - History
  - Models
  - User
- **Utility Packages**:
  - `ttsscript` - TTS script authoring
  - `voices` - Voice constants and metadata
- Functional options pattern for client configuration
- Automatic API key from environment variable
- Comprehensive error handling
- Multiple output format support (MP3, PCM, μ-law)

[0.3.0]: https://github.com/agentplexus/go-elevenlabs/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/agentplexus/go-elevenlabs/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/agentplexus/go-elevenlabs/releases/tag/v0.1.0
