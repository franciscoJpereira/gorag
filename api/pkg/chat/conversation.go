package chat

import apiinterface "ragAPI/pkg/apiInterface"

// Keeps the chat history
type Chat struct {
	messages []apiinterface.ChatMessage
}

func NewChat(messages ...apiinterface.ChatMessage) *Chat {
	return &Chat{
		messages,
	}
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
