package chat

import (
	"fmt"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
)

// Keeps the chat history
type Chat struct {
	messages []apiinterface.ChatMessage
}

func NewChat(messages ...apiinterface.ChatMessage) *Chat {
	return &Chat{
		messages,
	}
}

func (c *Chat) Store(chatStore store.Store, name string) error {
	return chatStore.Store(
		store.ChatHistory{
			ChatName: name,
			Messages: c.messages,
		},
	)
}

func (c *Chat) update(
	userMessage string,
	response apiinterface.ChatMessage,
) {
	userChatMessage := apiinterface.ChatMessage{
		Role:    "User",
		Content: userMessage,
	}
	c.messages = append(c.messages, userChatMessage, response)
}

func (c *Chat) GetHistory() []apiinterface.ChatMessage {
	return c.messages
}

func (c *Chat) NewMessageWithContext(
	message string,
	chatClient apiinterface.ApiInterface,
	context []string,
) (response apiinterface.ChatMessage, err error) {
	message = fmt.Sprintf("Answer the following: %s\nConsidering this information: ", message)
	for index, contextData := range context {
		message = fmt.Sprintf("%s\n\t%d. %s", message, index+1, contextData)
	}
	message = fmt.Sprintf("%s\nAnswer in the SAME language as the initial Query.", message)
	response, err = chatClient.Send(message, c.messages...)
	if err != nil {
		return
	}
	c.update(message, response)
	return
}

func (c *Chat) NewMessage(
	message string,
	chatClient apiinterface.ApiInterface,
) (response apiinterface.ChatMessage, err error) {

	response, err = chatClient.Send(message, c.messages...)

	c.update(
		message,
		response,
	)

	return
}
