package store

import apiinterface "ragAPI/pkg/apiInterface"

type Store interface {
	// Lists all chat names available
	ListChats() ([]string, error)
	//Gets a specific chat from the store
	Get(chatName string) (ChatHistory, error)
	//Stores the history into the store again
	Store(chat ChatHistory) error
}

type ChatHistory struct {
	ChatName string                     `json:"name"`
	Messages []apiinterface.ChatMessage `json:"history"`
}
