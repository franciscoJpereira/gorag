package apiinterface

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiInterface interface {
	Models() ([]string, error)
	Send(message string, thread ...ChatMessage) (ChatMessage, error)
}

type BasicInterface struct{}

func (b *BasicInterface) Models() ([]string, error) {
	return []string{"basic"}, nil
}

func (b *BasicInterface) Send(message string, thread ...ChatMessage) (ChatMessage, error) {
	return ChatMessage{
		Role:    "system",
		Content: "response",
	}, nil
}
