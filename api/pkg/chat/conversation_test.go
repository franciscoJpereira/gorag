package chat_test

import (
	"fmt"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat"
	"testing"
)

type MockedInterface struct{}

func (s *MockedInterface) Models() (models []string, err error) { return }

func (s *MockedInterface) Send(
	userMessage string,
	chat ...apiinterface.ChatMessage,
) (response apiinterface.ChatMessage, err error) {
	return apiinterface.ChatMessage{
		Role:    "System",
		Content: fmt.Sprintf("Response to: %s", userMessage),
	}, nil
}

func TestConversationResponse(t *testing.T) {
	client := &MockedInterface{}
	chat := chat.NewChat()

	response, err := chat.NewMessage("Hello, Bot!", client)

	if err != nil {
		t.Fatalf("Failed with error: %s", err)
	}

	if response.Role != "System" {
		t.Fatalf("Response Role invalid (System != %s)\n", response.Role)
	}

	if response.Content != "Response to: Hello, Bot!" {
		t.Fatalf("Response invalid, got: %s\n", response.Content)
	}
}

func TestHistoryResponse(t *testing.T) {
	client := &MockedInterface{}
	chat := chat.NewChat()

	_, _ = chat.NewMessage("Hello, Bot!", client)

	history := chat.GetHistory()

	if len(history) != 2 {
		t.Fatalf("Invalid history retrieved with length: %d", len(history))
	}

	if history[0].Role != "User" && history[0].Content != "Hello, Bot!" {
		t.Fatalf("First history message invalid: %v", history[0])
	}

	if history[1].Role != "Sytem" && history[1].Content != "Response to: Hello, Bot!" {
		t.Fatalf("Second history message invalid: %v", history[1])
	}
}
