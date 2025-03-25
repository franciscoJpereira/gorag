package localnet_test

import (
	"fmt"
	"os"
	"path/filepath"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	knowledgebase "ragAPI/pkg/knowledge-base"
	localnet "ragAPI/pkg/local-net"
	"testing"
)

const (
	TEST_STORE_PATH = "../../test/"
)

func createMockedRAG() *pkg.RAG {
	jsonStore, err := store.NewJsonStore(TEST_STORE_PATH)
	if err != nil {
		panic(fmt.Sprintf("Error Creating JSON store: %s\n", err))
	}
	basicInterface := &apiinterface.BasicInterface{}
	kb := knowledgebase.NewBasicBase()

	return &pkg.RAG{
		Api:       basicInterface,
		ChatStore: jsonStore,
		Kb:        kb,
	}
}

func removeFile(filename string) {
	absPath, _ := filepath.Abs(filename)
	if err := os.Remove(absPath); err != nil {
		panic(err)
	}
}

func TestAvailableKB(t *testing.T) {
	rag := createMockedRAG()
	rag.CreateKB("KB1")
	rag.CreateKB("KB2")
	controler := localnet.NewLocalControler(rag)
	kbs, err := controler.GetAvailableKBs()
	if err != nil {
		t.Fatalf("Failed to load kbs: %s\n", kbs)
	}
	if len(kbs) != 2 {
		t.Fatalf("Different amount of kbs: %d\n", len(kbs))
	}
	for i := range 2 {
		if kbs[i] != fmt.Sprintf("KB%d", i+1) {
			t.Fatalf("Different kb as expected")
		}
	}
}

func TestCreateKB(t *testing.T) {
	rag := createMockedRAG()
	controler := localnet.NewLocalControler(rag)
	if err := controler.CreateKB("KB1"); err != nil {
		t.Fatalf("Failed creating KB: %s\n", err)
	}
	kbs, _ := controler.GetAvailableKBs()
	if len(kbs) != 1 {
		t.Fatalf("Different length as expected: %d\n", len(kbs))
	}
	if kbs[0] != "KB1" {
		t.Fatalf("Different KB name: %s\n", kbs[0])
	}
}

func TestAddDataToKB(t *testing.T) {
	rag := createMockedRAG()
	rag.CreateKB("KB1")
	controler := localnet.NewLocalControler(rag)
	err := controler.AddDataToKB(pkg.KBAddDataInstruct{
		Data:   []string{"Some data"},
		KBName: "KB1",
	})
	if err != nil {
		t.Fatalf("Failed adding data to KB: %s", err)
	}
	data := rag.Kb.Retrieve("KB1", "")
	if len(data) != 1 {
		t.Fatalf("Data was not stored in Kb: %d", len(data))
	}
	if data[0] != "Some data" {
		t.Fatalf("Incorrect data stored in kb: %s", data[0])
	}
}

func TestSingleSHotMessage(t *testing.T) {
	rag := createMockedRAG()
	controler := localnet.NewLocalControler(rag)
	response, err := controler.SingleShotMessage(
		pkg.MessageInstruct{
			Message: "Hola",
		},
	)
	if err != nil {
		t.Fatalf("Failed sending message: %s\n", err)
	}
	if len(response.Ctx) != 0 {
		t.Fatalf("Failed context length: %d\n", len(response.Ctx))
	}
	if response.Query != "Hola" {
		t.Fatalf("Failed response query: %s\n", response.Query)
	}
}

func TestMessageToChat(t *testing.T) {
	rag := createMockedRAG()
	rag.NewChatMessage(pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: "Hola",
		},
		NewChat:  true,
		ChatName: pkg.EncodeBase64("Chat1"),
	})
	controler := localnet.NewLocalControler(rag)
	r, err := controler.SendNewMessageToChat(
		pkg.ChatInstruct{
			Message: pkg.MessageInstruct{
				Message: "Hola 2",
			},
			ChatName: "Chat1",
		},
	)
	if err != nil {
		t.Fatalf("Failed sending message to chat: %s\n", err)
	}
	if r.Query != "Hola 2" {
		t.Fatalf("Query different from expected: %s\n", r.Query)
	}
	removeFile("../../test/" + pkg.EncodeBase64("Chat1") + ".json")
}

func TestRetrievingChatNames(t *testing.T) {
	rag := createMockedRAG()
	chatNames := []string{
		"Chat1",
		"Chat2",
		"Chat3",
		"Chat4",
	}
	for _, chatname := range chatNames {
		rag.NewChatMessage(
			pkg.ChatInstruct{
				Message: pkg.MessageInstruct{
					Message: "Hola",
				},
				NewChat:  true,
				ChatName: pkg.EncodeBase64(chatname),
			},
		)
	}
	controler := localnet.NewLocalControler(rag)
	chats, err := controler.RetrieveAvailableChats()
	if err != nil {
		t.Fatalf("Failed retrieving chat names: %s\n", err)
	}
	for index, chatname := range chatNames {
		if chats[index] != chatname {
			t.Fatalf("Failed chat name at %d: %s\n", index, chats[index])
		}
		removeFile(fmt.Sprintf("../../test/%s.json", pkg.EncodeBase64(chatname)))
	}
}

func TestRetrieveChat(t *testing.T) {
	rag := createMockedRAG()
	chatname := "Chat1"
	rag.NewChatMessage(pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: "Created",
		},
		NewChat:  true,
		ChatName: pkg.EncodeBase64(chatname),
	})
	controler := localnet.NewLocalControler(rag)
	ch, err := controler.RetrieveChat(chatname)
	if err != nil {
		t.Fatalf("Failed to retrieve Chat: %s\n", err)
	}
	if ch.ChatName != chatname {
		t.Fatalf("Invalid chatname: %s\n", ch.ChatName)
	}
	if len(ch.Messages) != 2 {
		t.Fatalf("Invalid chat messages: %v\n", ch.Messages)
	}
	if ch.Messages[0].Content != "Created" {
		t.Fatalf("Invalid message content: %s\n", ch.Messages[0].Content)
	}
	removeFile("../../test/" + pkg.EncodeBase64("Chat1") + ".json")
}
