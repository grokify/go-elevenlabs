package elevenlabs

import (
	"context"
	"testing"
)

func TestHistoryList_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.History().List(context.Background(), nil)
	if err != nil {
		t.Fatalf("History().List() error = %v", err)
	}
	if resp == nil {
		t.Fatal("History().List() returned nil")
	}

	// Items can be empty if user has no history
	t.Logf("History items: %d, HasMore: %v", len(resp.Items), resp.HasMore)
}

func TestHistoryListWithOptions_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	opts := &HistoryListOptions{
		PageSize: 5,
	}

	resp, err := client.History().List(context.Background(), opts)
	if err != nil {
		t.Fatalf("History().List() error = %v", err)
	}
	if resp == nil {
		t.Fatal("History().List() returned nil")
	}

	// Verify page size is respected (if there are items)
	if len(resp.Items) > 5 {
		t.Errorf("History().List() returned %d items, expected max 5", len(resp.Items))
	}
}

func TestHistoryGetValidation(t *testing.T) {
	client, _ := NewClient()

	_, err := client.History().Get(context.Background(), "")
	if err == nil {
		t.Error("Get('') should return error")
	}

	_, err = client.History().GetAudio(context.Background(), "")
	if err == nil {
		t.Error("GetAudio('') should return error")
	}

	err = client.History().Delete(context.Background(), "")
	if err == nil {
		t.Error("Delete('') should return error")
	}
}
