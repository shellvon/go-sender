package core

import (
	"context"
	"testing"
	"time"
)

// MockProvider implements Provider for testing
type MockProvider struct {
	name string
}

func (m *MockProvider) Send(ctx context.Context, message Message) error {
	return nil
}

func (m *MockProvider) Name() string {
	return m.name
}

// MockMessage implements Message for testing
type MockMessage struct {
	DefaultMessage
	content string
}

func (m *MockMessage) Validate() error {
	return nil
}

func (m *MockMessage) ProviderType() ProviderType {
	return ProviderTypeSMS
}

func (m *MockMessage) Content() string {
	return m.content
}

func TestContextSerializationRoundTrip(t *testing.T) {
	// Test the serialization and deserialization of context information
	originalCtx := context.Background()
	originalCtx = WithCtxItemName(originalCtx, "test-bot")

	// Create metadata
	metadata := make(map[string]interface{})

	// Serialize
	opts := &SendOptions{
		Priority: 10,
		Timeout:  time.Second * 60,
	}

	serializedMetadata, err := serializeSendOptions(originalCtx, opts, metadata)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Deserialize
	restoredCtx, restoredOpts, err := deserializeSendOptions(context.Background(), serializedMetadata)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Verify context information was restored
	if itemName := GetItemNameFromCtx(restoredCtx); itemName != "test-bot" {
		t.Errorf("Expected item name 'test-bot', got '%s'", itemName)
	}

	// Verify SendOptions were restored
	if restoredOpts.Priority != 10 {
		t.Errorf("Expected priority 10, got %d", restoredOpts.Priority)
	}

	if restoredOpts.Timeout != time.Second*60 {
		t.Errorf("Expected timeout 60s, got %v", restoredOpts.Timeout)
	}
}

func TestContextSerializationWithEmptyContext(t *testing.T) {
	// Test serialization with empty context
	ctx := context.Background()
	metadata := make(map[string]interface{})
	opts := &SendOptions{
		Priority: 5,
	}

	serializedMetadata, err := serializeSendOptions(ctx, opts, metadata)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Deserialize
	restoredCtx, restoredOpts, err := deserializeSendOptions(context.Background(), serializedMetadata)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Verify no context information was added
	if itemName := GetItemNameFromCtx(restoredCtx); itemName != "" {
		t.Errorf("Expected empty item name, got '%s'", itemName)
	}

	// Verify SendOptions were restored
	if restoredOpts.Priority != 5 {
		t.Errorf("Expected priority 5, got %d", restoredOpts.Priority)
	}
}
