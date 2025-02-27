package apiinterface

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiInterface interface {
	Models() ([]string, error)
	Send(message string, thread ...ChatMessage) (ChatMessage, error)
}
