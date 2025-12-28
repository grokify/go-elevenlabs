package elevenlabs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketSTTService handles real-time speech-to-text via WebSocket.
type WebSocketSTTService struct {
	client *Client
}

// WebSocketSTTOptions configures the WebSocket STT connection.
type WebSocketSTTOptions struct {
	// ModelID is the transcription model to use.
	// Options: "scribe_v1" (default), "scribe_v1_experimental"
	ModelID string

	// LanguageCode is the expected language (e.g., "en", "es").
	// If not specified, language will be auto-detected.
	LanguageCode string

	// SampleRate is the audio sample rate in Hz.
	// Common values: 8000, 16000, 22050, 44100
	SampleRate int

	// Encoding is the audio encoding format.
	// Options: "pcm_s16le" (default), "pcm_mulaw"
	Encoding string

	// EnablePartials enables partial/interim transcription results.
	EnablePartials bool

	// EnableWordTimestamps enables word-level timing information.
	EnableWordTimestamps bool

	// MaxAlternatives is the maximum number of transcription alternatives.
	MaxAlternatives int
}

// DefaultWebSocketSTTOptions returns default options for real-time STT.
func DefaultWebSocketSTTOptions() *WebSocketSTTOptions {
	return &WebSocketSTTOptions{
		ModelID:              "scribe_v1",
		SampleRate:           16000,
		Encoding:             "pcm_s16le",
		EnablePartials:       true,
		EnableWordTimestamps: true,
	}
}

// WebSocketSTTConnection represents an active WebSocket STT connection.
type WebSocketSTTConnection struct {
	conn    *websocket.Conn
	options *WebSocketSTTOptions
	mu      sync.Mutex
	closed  bool

	// Channels for async operation
	transcriptOut chan *STTTranscript
	errChan       chan error
	closeChan     chan struct{}
	closeOnce     sync.Once
}

// STTTranscript represents a transcription result.
type STTTranscript struct {
	// Text is the transcribed text.
	Text string `json:"text"`

	// IsFinal indicates if this is a final (non-partial) result.
	IsFinal bool `json:"is_final"`

	// Confidence is the confidence score (0.0 to 1.0).
	Confidence float64 `json:"confidence,omitempty"`

	// Words contains word-level timing if enabled.
	Words []STTWord `json:"words,omitempty"`

	// LanguageCode is the detected language.
	LanguageCode string `json:"language_code,omitempty"`

	// StartTime is the start time in seconds.
	StartTime float64 `json:"start_time,omitempty"`

	// EndTime is the end time in seconds.
	EndTime float64 `json:"end_time,omitempty"`
}

// STTWord represents a single word with timing.
type STTWord struct {
	Word       string  `json:"word"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Confidence float64 `json:"confidence,omitempty"`
}

// sttWSInitMessage is the initial configuration message.
type sttWSInitMessage struct {
	Type                 string `json:"type"`
	SampleRate           int    `json:"sample_rate,omitempty"`
	Encoding             string `json:"encoding,omitempty"`
	LanguageCode         string `json:"language_code,omitempty"`
	EnablePartials       bool   `json:"enable_partials,omitempty"`
	EnableWordTimestamps bool   `json:"enable_word_timestamps,omitempty"`
	MaxAlternatives      int    `json:"max_alternatives,omitempty"`
}

// sttWSAudioMessage is an audio data message.
type sttWSAudioMessage struct {
	Type  string `json:"type"`
	Audio string `json:"audio"` // Base64 encoded audio
}

// sttWSControlMessage is a control message.
type sttWSControlMessage struct {
	Type string `json:"type"`
}

// sttWSResponse is the WebSocket response from STT.
type sttWSResponse struct {
	Type         string    `json:"type"`
	Text         string    `json:"text,omitempty"`
	IsFinal      bool      `json:"is_final,omitempty"`
	Confidence   float64   `json:"confidence,omitempty"`
	Words        []STTWord `json:"words,omitempty"`
	LanguageCode string    `json:"language_code,omitempty"`
	StartTime    float64   `json:"start_time,omitempty"`
	EndTime      float64   `json:"end_time,omitempty"`
	Error        string    `json:"error,omitempty"`
	Message      string    `json:"message,omitempty"`
}

// Connect establishes a WebSocket connection for real-time STT.
func (s *WebSocketSTTService) Connect(ctx context.Context, opts *WebSocketSTTOptions) (*WebSocketSTTConnection, error) {
	if opts == nil {
		opts = DefaultWebSocketSTTOptions()
	}

	// Build WebSocket URL
	wsURL, err := s.buildWebSocketURL(opts)
	if err != nil {
		return nil, err
	}

	// Create dialer with context
	dialer := websocket.Dialer{
		HandshakeTimeout: 0, // Use context timeout
	}

	// Add headers
	headers := http.Header{}
	headers.Set("xi-api-key", s.client.apiKey)

	// Connect
	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	wsc := &WebSocketSTTConnection{
		conn:          conn,
		options:       opts,
		transcriptOut: make(chan *STTTranscript, 100),
		errChan:       make(chan error, 1),
		closeChan:     make(chan struct{}),
	}

	// Send initial configuration
	if err := wsc.sendInit(); err != nil {
		conn.Close()
		return nil, err
	}

	// Start reading responses
	go wsc.readLoop()

	return wsc, nil
}

func (s *WebSocketSTTService) buildWebSocketURL(opts *WebSocketSTTOptions) (string, error) {
	baseURL := s.client.baseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Convert HTTP URL to WebSocket URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = "/v1/speech-to-text/realtime"

	// Add query parameters
	q := u.Query()
	if opts.ModelID != "" {
		q.Set("model_id", opts.ModelID)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (wsc *WebSocketSTTConnection) sendInit() error {
	msg := sttWSInitMessage{
		Type:                 "config",
		SampleRate:           wsc.options.SampleRate,
		Encoding:             wsc.options.Encoding,
		EnablePartials:       wsc.options.EnablePartials,
		EnableWordTimestamps: wsc.options.EnableWordTimestamps,
	}

	if wsc.options.LanguageCode != "" {
		msg.LanguageCode = wsc.options.LanguageCode
	}

	if wsc.options.MaxAlternatives > 0 {
		msg.MaxAlternatives = wsc.options.MaxAlternatives
	}

	return wsc.sendJSON(msg)
}

func (wsc *WebSocketSTTConnection) sendJSON(msg any) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return fmt.Errorf("connection closed")
	}

	return wsc.conn.WriteJSON(msg)
}

func (wsc *WebSocketSTTConnection) readLoop() {
	defer wsc.closeChannels()

	for {
		select {
		case <-wsc.closeChan:
			return
		default:
		}

		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				select {
				case wsc.errChan <- err:
				default:
				}
			}
			return
		}

		var resp sttWSResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			select {
			case wsc.errChan <- fmt.Errorf("failed to parse response: %w", err):
			default:
			}
			continue
		}

		// Check for errors
		if resp.Error != "" || (resp.Type == "error" && resp.Message != "") {
			errMsg := resp.Error
			if errMsg == "" {
				errMsg = resp.Message
			}
			select {
			case wsc.errChan <- fmt.Errorf("server error: %s", errMsg):
			default:
			}
			continue
		}

		// Handle transcript responses
		if resp.Type == "transcript" || resp.Text != "" {
			transcript := &STTTranscript{
				Text:         resp.Text,
				IsFinal:      resp.IsFinal,
				Confidence:   resp.Confidence,
				Words:        resp.Words,
				LanguageCode: resp.LanguageCode,
				StartTime:    resp.StartTime,
				EndTime:      resp.EndTime,
			}
			select {
			case wsc.transcriptOut <- transcript:
			case <-wsc.closeChan:
				return
			}
		}
	}
}

func (wsc *WebSocketSTTConnection) closeChannels() {
	wsc.closeOnce.Do(func() {
		close(wsc.closeChan)
		close(wsc.transcriptOut)
	})
}

// SendAudio sends audio data for transcription.
// The audio should be in the format specified in WebSocketSTTOptions.
func (wsc *WebSocketSTTConnection) SendAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	msg := sttWSAudioMessage{
		Type:  "audio",
		Audio: base64.StdEncoding.EncodeToString(audio),
	}

	return wsc.sendJSON(msg)
}

// EndStream signals that no more audio will be sent.
// This allows the server to finalize any pending transcription.
func (wsc *WebSocketSTTConnection) EndStream() error {
	msg := sttWSControlMessage{
		Type: "end_of_stream",
	}
	return wsc.sendJSON(msg)
}

// Transcripts returns a channel that receives transcription results.
func (wsc *WebSocketSTTConnection) Transcripts() <-chan *STTTranscript {
	return wsc.transcriptOut
}

// Errors returns a channel that receives errors from the connection.
func (wsc *WebSocketSTTConnection) Errors() <-chan error {
	return wsc.errChan
}

// Close closes the WebSocket connection gracefully.
func (wsc *WebSocketSTTConnection) Close() error {
	wsc.mu.Lock()
	if wsc.closed {
		wsc.mu.Unlock()
		return nil
	}
	wsc.closed = true
	wsc.mu.Unlock()

	// Send end of stream
	_ = wsc.EndStream()

	// Close the connection
	wsc.closeChannels()
	return wsc.conn.Close()
}

// StreamAudio is a convenience method that streams audio from a channel.
// It handles ending the stream automatically when the input channel closes.
func (wsc *WebSocketSTTConnection) StreamAudio(ctx context.Context, audioStream <-chan []byte) (<-chan *STTTranscript, <-chan error) {
	transcriptOut := make(chan *STTTranscript, 100)
	errOut := make(chan error, 1)

	go func() {
		defer close(transcriptOut)
		defer close(errOut)

		// Forward transcripts from connection
		done := make(chan struct{})
		go func() {
			defer close(done)
			for transcript := range wsc.Transcripts() {
				select {
				case transcriptOut <- transcript:
				case <-ctx.Done():
					return
				}
			}
		}()

		// Send audio as it arrives
		for {
			select {
			case audio, ok := <-audioStream:
				if !ok {
					// Input stream closed, end stream and wait for remaining transcripts
					if err := wsc.EndStream(); err != nil {
						errOut <- err
						return
					}
					<-done
					return
				}
				if err := wsc.SendAudio(audio); err != nil {
					errOut <- err
					return
				}
			case err := <-wsc.Errors():
				errOut <- err
				return
			case <-ctx.Done():
				errOut <- ctx.Err()
				return
			}
		}
	}()

	return transcriptOut, errOut
}
