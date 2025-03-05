package pkg

// Data to add data to a Knowledge base
type KBAddDataInstruct struct {
	Data   []string `json:"data"`
	KBName string   `json:"KBName"`
	Create bool     `json:"Create"`
}

// Data to send a Message
type MessageInstruct struct {
	Message string `json:"Message"`
	UseKB   bool   `json:"KB"`
	KBName  string `json:"KBName"`
}

// Data to send a new message to a chat
type ChatInstruct struct {
	Message  MessageInstruct `json:"Message"`
	NewChat  bool            `json:"NewChat"`
	ChatName string          `json:"ChatName"`
}

// Message Response
type MessageResponse struct {
	Ctx      []string `json:"Context"`
	Query    string   `json:"Query"`
	Response string   `json:"Response"`
}
