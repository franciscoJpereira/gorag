package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	knowledgebase "ragAPI/pkg/knowledge-base"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

const (
	TEST_STORE_PATH = "../test/"
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

func createChat(chatname string, rag *pkg.RAG) {
	data := pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: "Created",
		},
		NewChat:  true,
		ChatName: chatname,
	}
	if _, err := rag.NewChatMessage(data); err != nil {
		panic(err.Error())
	}
}

func TestCreateKB(t *testing.T) {
	e := echo.New()
	rag := createMockedRAG()
	req := httptest.NewRequest(http.MethodPost, "/knowledge-base/newKB", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("KBName")
	c.SetParamValues("newKB")
	c.Set(RAGKey, rag)
	if err := CreateKB(c); err != nil {
		t.Fatalf("Failed witn error: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Invalid status: %d\n", rec.Code)
	}
	//Test trying to create the same collection once more
	req = httptest.NewRequest(http.MethodPost, "/knowledge-base/newKB", strings.NewReader(""))
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("KBName")
	c.SetParamValues("newKB")
	c.Set(RAGKey, rag)
	if err := CreateKB(c); err != nil {
		t.Fatalf("Double call created an error: %s\n", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Invalid status on double request: %d\n", rec.Code)
	}
}

func TestGetAvailableKBs(t *testing.T) {
	e := echo.New()
	rag := createMockedRAG()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := GetAvailableKBs(c); err != nil {
		t.Fatalf("First request failed with err: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("INvalid response code: %d\n", rec.Code)
	}
	var response []string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not unmarshal response: %s\n", err)
	}
	if len(response) != 0 {
		t.Fatalf("Invalid response: %v\n", response)
	}
	rag.CreateKB("KB1")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := GetAvailableKBs(c); err != nil {
		t.Fatalf("Second request failed with err: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Invalid response code: %d\n", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not unmarshal response: %s\n", err)
	}
	if len(response) != 1 {
		t.Fatalf("Invalid response: %v\n", response)
	}
	if response[0] != "KB1" {
		t.Fatalf("Invalid KB name: %s\n", response[0])
	}
}

func TestAddDataToKB(t *testing.T) {
	rag := createMockedRAG()
	data := pkg.KBAddDataInstruct{
		Create: true,
		KBName: "TestKB",
		Data:   []string{"Some data"},
	}
	marshaled, _ := json.Marshal(data)
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewReader(marshaled),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := AddDataToKB(c); err != nil {
		t.Fatalf("Failed calling method: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed with status code: %d", rec.Code)
	}
	data.Create = false
	marshaled, _ = json.Marshal(data)
	req = httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewReader(marshaled),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := AddDataToKB(c); err != nil {
		t.Fatalf("Failed calling method 2: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed 2 with status code: %d\n", rec.Code)
	}
	recovered := rag.Kb.Retrieve(data.KBName, "")
	if len(recovered) != 2 {
		t.Fatalf("Failed to load data correctly: %v\n", recovered)
	}
}

func TestSingleShotMessageWithoutContext(t *testing.T) {
	rag := createMockedRAG()
	instruction := pkg.MessageInstruct{
		Message: "Message 1",
	}
	marshaled, _ := json.Marshal(instruction)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshaled))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := SingleShotMessage(c); err != nil {
		t.Fatalf("Failed message: %s\n", err)
	}
	var response pkg.MessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed unmarshaling response: %s\n", err)
	}
	if len(response.Ctx) != 0 {
		t.Fatalf("Context was not null: %v\n", response.Ctx)
	}
	if response.Query != "Message 1" {
		t.Fatalf("Invalid query: %s\n", response.Query)
	}
}

func TestSingleShotMessageWithContext(t *testing.T) {
	rag := createMockedRAG()
	e := echo.New()
	data := pkg.MessageInstruct{
		Message: "Message 1",
		UseKB:   true,
		KBName:  "KB1",
	}
	marshaled, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshaled))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	rag.CreateKB("KB1")
	rag.Kb.AddDataToCollection("KB1", []string{"Context data"})
	c.Set(RAGKey, rag)
	if err := SingleShotMessage(c); err != nil {
		t.Fatalf("Failed single shot message: %s\n", err)
	}
	var response pkg.MessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed unmarshaling JSON: %s\n", err)
	}
	if len(response.Ctx) != 1 {
		t.Fatalf("Invalid context length: %d\n", len(response.Ctx))
	}
	if response.Ctx[0] != "Context data" {
		t.Fatalf("Failed data: %s\n", response.Ctx[0])
	}
	if response.Query != "Message 1" {
		t.Fatalf("Failed query data: %s\n", response.Query)
	}
}

func TestCreateChat(t *testing.T) {
	rag := createMockedRAG()
	data := pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: "Message 1",
		},
		ChatName: "Chat 1",
		NewChat:  true,
	}
	marshaled, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshaled))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := SendNewMessageToChat(c); err != nil {
		t.Fatalf("Failed call to function: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed status code: %d\n", rec.Code)
	}
	var response pkg.MessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed unmarshaling json: %s\n", err)
	}
	if response.Query != data.Message.Message {
		t.Fatalf("Failed response query: %s\n", response.Query)
	}
	removeFile("../test/Chat 1.json")
}

func TestNewMessageToExistingChat(t *testing.T) {
	rag := createMockedRAG()
	data := pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: "Message 1",
		},
		ChatName: "Chat 1",
		NewChat:  true,
	}
	rag.NewChatMessage(data)
	data.Message.Message = "Message 2"
	data.NewChat = false
	marshaled, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshaled))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := SendNewMessageToChat(c); err != nil {
		t.Fatalf("Failed sending message to Chat: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed response code: %d\n", rec.Code)
	}

	removeFile("../test/Chat 1.json")
}

func TestNoChatRetrieval(t *testing.T) {
	rag := createMockedRAG()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := RetrieveAvailableChats(c); err != nil {
		t.Fatalf("Failed retrieving chat: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed status code: %d\n", rec.Code)
	}
	var response []string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed unmarshaling body: %s\n", err)
	}
	if len(response) != 0 {
		t.Fatalf("Failed response: %v\n", response)
	}
}

func TestMultipleChatRetrieval(t *testing.T) {
	rag := createMockedRAG()
	createChat("Chat 1", rag)
	createChat("Chat 2", rag)
	createChat("Chat 3", rag)
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set(RAGKey, rag)
	if err := RetrieveAvailableChats(c); err != nil {
		t.Fatalf("Failed retrieving chat: %s\n", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed status code: %d\n", rec.Code)
	}
	var response []string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed unmarshaling response: %s\n", err)
	}
	if len(response) != 3 {
		t.Fatalf("Failed response length: %d\n", len(response))
	}
	for i := range 3 {
		if response[i] != fmt.Sprintf("Chat %d", i+1) {
			t.Fatalf("Failed at response %d with value: %s", i, response[i])
		}
	}
	removeFile("../test/Chat 1.json")
	removeFile("../test/Chat 2.json")
	removeFile("../test/Chat 3.json")
}
