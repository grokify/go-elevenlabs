package elevenlabs

import (
	"context"
	"testing"
)

func TestUserGetInfo_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	user, err := client.User().GetInfo(context.Background())
	if err != nil {
		t.Fatalf("User().GetInfo() error = %v", err)
	}
	if user == nil {
		t.Fatal("User().GetInfo() returned nil")
	}
	if user.UserID == "" {
		t.Error("User.UserID is empty")
	}
	if user.Subscription == nil {
		t.Error("User.Subscription is nil")
	}
}

func TestUserGetSubscription_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	sub, err := client.User().GetSubscription(context.Background())
	if err != nil {
		t.Fatalf("User().GetSubscription() error = %v", err)
	}
	if sub == nil {
		t.Fatal("User().GetSubscription() returned nil")
	}
	if sub.Tier == "" {
		t.Error("Subscription.Tier is empty")
	}
}

func TestUserGetCharactersRemaining_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	remaining, err := client.User().GetCharactersRemaining(context.Background())
	if err != nil {
		t.Fatalf("User().GetCharactersRemaining() error = %v", err)
	}
	// remaining can be any value including 0, just check it doesn't error
	t.Logf("Characters remaining: %d", remaining)
}

func TestSubscriptionCharactersRemaining(t *testing.T) {
	tests := []struct {
		name     string
		sub      *Subscription
		expected int
	}{
		{
			name:     "normal case",
			sub:      &Subscription{CharacterCount: 100, CharacterLimit: 1000},
			expected: 900,
		},
		{
			name:     "over limit",
			sub:      &Subscription{CharacterCount: 1500, CharacterLimit: 1000},
			expected: 0,
		},
		{
			name:     "at limit",
			sub:      &Subscription{CharacterCount: 1000, CharacterLimit: 1000},
			expected: 0,
		},
		{
			name:     "no usage",
			sub:      &Subscription{CharacterCount: 0, CharacterLimit: 1000},
			expected: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sub.CharactersRemaining()
			if got != tt.expected {
				t.Errorf("CharactersRemaining() = %d, want %d", got, tt.expected)
			}
		})
	}
}
