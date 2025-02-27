package store_test

import (
	"fmt"
	"os"
	"path/filepath"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	"testing"
)

func TestListingChats(t *testing.T) {
	chatStore, err := filepath.Abs("../../../test/store")
	if err != nil {
		t.Fatalf("Filepath failed: %s", err)
	}
	store, err := store.NewJsonStore(chatStore)
	if err != nil {
		t.Fatalf("Failed creating store: %s", err)
	}

	chats, err := store.ListChats()
	if err != nil {
		t.Fatalf("Failed to list all chats: %s", err)
	}
	if len(chats) != 1 {
		t.Fatalf("Expected only one chat to be valid, got %d\n%s", len(chats), chats)
	}
	if chats[0] != "test1" {
		t.Fatalf("Got %s instead of test1 as first value", chats[0])
	}
}

func TestGetChat(t *testing.T) {
	chatStore, err := filepath.Abs("../../../test/store")
	if err != nil {
		t.Fatalf("Filepath failed: %s", err)
	}
	store, err := store.NewJsonStore(chatStore)
	if err != nil {
		t.Fatalf("Failed creating store: %s", err)
	}
	chat, err := store.Get("test1")
	if err != nil {
		t.Fatalf("Failed geting chat test1: %s", err)
	}
	if chat.ChatName != "test1" {
		t.Fatalf("Chat name invalid, got: %s", chat.ChatName)
	}
	if len(chat.Messages) != 2 {
		t.Fatalf("Chat len invalid, got: %d", len(chat.Messages))
	}
	for _, message := range chat.Messages {
		if message.Content != fmt.Sprintf("Hello, %s", message.Role) {
			t.Fatalf("Chat message invalid, got: %s", message.Content)
		}
	}
}

func TestStoreChat(t *testing.T) {
	chatStore, err := filepath.Abs("../../../test/store")
	if err != nil {
		t.Fatalf("Filepath failed: %s", err)
	}
	jsonStore, err := store.NewJsonStore(chatStore)
	if err != nil {
		t.Fatalf("Failed creating store: %s", err)
	}
	history := store.ChatHistory{
		ChatName: "testing",
		Messages: []apiinterface.ChatMessage{
			{
				Role:    "Testing",
				Content: "Hello, Testing",
			},
		},
	}
	if err = jsonStore.Store(history); err != nil {
		t.Fatalf("Failed storing data: %s", err)
	}
	historyRecovered, err := jsonStore.Get(history.ChatName)
	if err != nil {
		t.Fatalf("Failed recovering chat: %s", err)
	}
	if historyRecovered.ChatName != history.ChatName {
		t.Fatalf("Different ChatNames: %s vs %s", history.ChatName, historyRecovered.ChatName)
	}
	if len(history.Messages) != len(historyRecovered.Messages) {
		t.Fatalf("Different Messages sizes")
	}
	for index, message := range history.Messages {
		recovered := historyRecovered.Messages[index]
		if message.Role != recovered.Role {
			t.Fatalf("Different role at %d: %s vs %s", index, message.Role, recovered.Role)
		}
		if message.Content != recovered.Content {
			t.Fatalf("Different Content at %d: %s vs %s", index, message.Content, recovered.Content)
		}
	}
	//Cleanup
	if err = os.Remove(fmt.Sprintf("%s/testing.json", chatStore)); err != nil {
		t.Logf("FAILED CLEAN-UP: %s", err)
	}
}
