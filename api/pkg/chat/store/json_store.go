package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// FS chat store
// Stores all chats on different json files (Where filename = chatname)
// The store path makes reference to the directory where the files are located.
type JsonChatStore struct {
	storePath string
}

func NewJsonStore(pathStore string) (*JsonChatStore, error) {

	if _, err := os.ReadDir(pathStore); err != nil {
		return nil, err
	}
	return &JsonChatStore{pathStore}, nil
}

func (s *JsonChatStore) ListChats() ([]string, error) {
	entries, err := os.ReadDir(s.storePath)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, entry := range entries {
		entryName := entry.Name()
		if strings.HasSuffix(entryName, ".json") {
			names = append(names, strings.TrimSuffix(entryName, ".json"))
		}
	}
	return names, nil
}

func (s *JsonChatStore) Get(chatName string) (history ChatHistory, err error) {
	var chat []byte
	chatPath := fmt.Sprintf("%s/%s.json", s.storePath, chatName)

	chat, err = os.ReadFile(chatPath)
	if err != nil {
		return
	}
	if err = json.Unmarshal(chat, &history); err != nil {
		return
	}
	return
}

func (s *JsonChatStore) Store(chat ChatHistory) (err error) {
	chatPath := fmt.Sprintf("%s/%s.json", s.storePath, chat.ChatName)
	data, err := json.Marshal(chat)

	if err != nil {
		return
	}

	err = os.WriteFile(chatPath, data, 0666)

	return
}
