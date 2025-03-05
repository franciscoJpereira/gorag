package pkg

import (
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat"
	"ragAPI/pkg/chat/store"
	knowledgebase "ragAPI/pkg/knowledge-base"
)

type ActiveChat struct {
	Name string
	Chat chat.Chat
}

type RAG struct {
	Api       apiinterface.ApiInterface
	ChatStore store.Store
	Kb        knowledgebase.BaseInterface
}

func (r *RAG) ListChats() ([]string, error) {
	return r.ChatStore.List()
}

func (r *RAG) RetrieveChat(chatName string) (store.ChatHistory, error) {
	return r.ChatStore.Get(chatName)
}

func (r *RAG) CreateKB(KBName string) error {
	return r.Kb.CreateColletion(KBName)
}

func (r *RAG) AddDataToKB(instruction KBAddDataInstruct) error {
	if instruction.Create {
		if err := r.CreateKB(instruction.KBName); err != nil {
			return err
		}
	}
	return r.Kb.AddDataToCollection(instruction.KBName, instruction.Data)
}

// Utility to send a new message
func (r *RAG) sendToChat(currentChat *chat.Chat, message MessageInstruct) (response MessageResponse, err error) {
	var responseText apiinterface.ChatMessage
	response.Query = message.Message
	if message.UseKB {
		context := r.Kb.Retrieve(message.KBName, message.Message)
		response.Ctx = context
		responseText, err = currentChat.NewMessageWithContext(message.Message, r.Api, context)
	} else {
		responseText, err = currentChat.NewMessage(message.Message, r.Api)
	}
	response.Response = responseText.Content
	return
}

// Message that will only use what's given by the store and
// the query to produce an answer
func (r *RAG) SingleShotMessage(message MessageInstruct) (response MessageResponse, err error) {
	currentChat := chat.NewChat()
	response, err = r.sendToChat(currentChat, message)
	return
}

// Message that will be added to a chat
func (r *RAG) NewChatMessage(message ChatInstruct) (response MessageResponse, err error) {
	var chatMessages store.ChatHistory
	if !message.NewChat {
		chatMessages, err = r.ChatStore.Get(message.ChatName)
	}
	currentChat := chat.NewChat(chatMessages.Messages...)
	if err != nil {
		return
	}
	response, err = r.sendToChat(currentChat, message.Message)
	if err != nil {
		return
	}
	err = currentChat.Store(r.ChatStore, message.ChatName)
	return
}
