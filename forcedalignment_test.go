package elevenlabs

import (
	"context"
	"strings"
	"testing"
)

func TestForcedAlignmentRequestValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *ForcedAlignmentRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil file",
			req:     &ForcedAlignmentRequest{Text: "test"},
			wantErr: true,
			errMsg:  "file",
		},
		{
			name:    "empty text",
			req:     &ForcedAlignmentRequest{File: strings.NewReader("audio"), Text: ""},
			wantErr: true,
			errMsg:  "text",
		},
		{
			name: "valid request",
			req: &ForcedAlignmentRequest{
				File:     strings.NewReader("audio"),
				Filename: "test.mp3",
				Text:     "Hello world",
			},
			wantErr: false, // Will fail at API level but pass validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.ForcedAlignment().Align(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("Align() should return error")
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

func TestForcedAlignmentService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.ForcedAlignment() == nil {
		t.Error("ForcedAlignment() returned nil")
	}
}

func TestAlignmentResponse(t *testing.T) {
	// Test response struct initialization
	resp := &ForcedAlignmentResponse{
		Loss: 0.1,
		Words: []AlignmentWord{
			{Text: "Hello", Start: 0.0, End: 0.5, Loss: 0.05},
			{Text: "world", Start: 0.6, End: 1.0, Loss: 0.08},
		},
		Characters: []AlignmentCharacter{
			{Text: "H", Start: 0.0, End: 0.1},
		},
	}

	if resp.Loss != 0.1 {
		t.Errorf("Loss = %f, want 0.1", resp.Loss)
	}
	if len(resp.Words) != 2 {
		t.Errorf("Words count = %d, want 2", len(resp.Words))
	}
	if len(resp.Characters) != 1 {
		t.Errorf("Characters count = %d, want 1", len(resp.Characters))
	}
}
