package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	RAGKey = "Rag"
)

// @Summary Get available Knowledge Bases
// @Description Returns a list with the names of all available knwoledge bases
// @Tags knowledge-base
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /knowledge-base [get]
func GetAvailableKBs(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG us not set")
	}
	result, err := rag.ListKBs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

// @Summary Create a new knowledge base
// @Description Create a new knowledge base
// @Tags knowledge-base
// @Param KBName path string true "KBName"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "KB already exists"
// @Router /knowledge-base/{KBName} [post]
func CreateKB(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	collectionName := c.Param("KBName")
	err := rag.Kb.CreateCollection(collectionName)
	if err != nil {
		return c.String(http.StatusBadRequest, "KB already exists")
	}
	return c.NoContent(http.StatusOK)
}

// @Summary Add data to a knowledge base
// @Description Add string data to a knwoledge base, it creates the KB if the flag is set
// @Tags knowledge-base
// @Param request body pkg.KBAddDataInstruct true "Data to add to the KB"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "KB does not exist"
// @Router /knowledge-base [post]
func AddDataToKB(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	var instruction KBAddDataInstruct
	if err := c.Bind(&instruction); err != nil {
		return c.String(http.StatusBadRequest, "Invalid data")
	}
	err := rag.AddDataToKB(instruction)
	if err != nil {
		return c.String(http.StatusBadRequest, "KB does not exist")
	}
	return c.NoContent(http.StatusOK)
}

// @Summary Send a one-shot message
// @Description Send a one-shot message to get a response
// @Tags chat
// @Param request body pkg.MessageInstruct true "Message to send"
// @Accept json
// @Produce json
// @Success 200 {object} pkg.MessageResponse
// @Failure 400 {string} string "Error sending message"
// @Router /message [post]
func SingleShotMessage(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	var message MessageInstruct
	if err := c.Bind(&message); err != nil {
		return c.String(http.StatusBadRequest, "Invalid data")
	}
	response, err := rag.SingleShotMessage(message)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending message")
	}
	return c.JSON(http.StatusOK, response)
}

// @Summary Send a new message to a chat
// @Description Sends a message to a chat and creates it if needed
// @Tags chat
// @Param request body pkg.ChatInstruct true "Message to send to Chat"
// @Accept json
// @Produce json
// @Success 200 {object} pkg.MessageResponse
// @Failure 400 {string} string "Error sending message"
// @Router /chat [post]
func SendNewMessageToChat(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	var message ChatInstruct
	if err := c.Bind(&message); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	response, err := rag.NewChatMessage(message)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

// @Summary Retrieves all available Chats
// @Description Returns the names of all existing chats
// @Tags chat
// @Param chatID query string false "Chat ID"
// @Accept json
// @produce json
// @Success 200 {array} string
// @Success 200 {object} store.ChatHistory
// @Router /chat [get]
func RetrieveAvailableChats(c echo.Context) error {
	rag, ok := c.Get(RAGKey).(*RAG)
	query := c.QueryParam("chatID")
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	if query != "" {
		chatHistory, err := rag.RetrieveChat(query)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, chatHistory)
	}
	names, err := rag.ListChats()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, names)
}
