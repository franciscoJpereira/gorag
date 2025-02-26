package apiinterface

type ChatMessage struct {
	Role    string
	Content string
}

type ApiInterface interface {
	Models() (string, error)
	Send(message string, thread ...ChatMessage) (ChatMessage, error)
}
