package main

import (
	"net/http"
	"ragAPI/pkg"

	"github.com/labstack/echo/v4"
)

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
	rag, ok := c.Get(RAGKey).(*pkg.RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	collectionName := c.Param("KBName")
	err := rag.Kb.CreateColletion(collectionName)
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
	rag, ok := c.Get(RAGKey).(*pkg.RAG)
	if !ok {
		return c.String(http.StatusInternalServerError, "RAG is not set")
	}
	var instruction pkg.KBAddDataInstruct
	if err := c.Bind(&instruction); err != nil {
		return c.String(http.StatusBadRequest, "Invalid data")
	}
	err := rag.AddDataToKB(instruction)
	if err != nil {
		return c.String(http.StatusBadRequest, "KB does not exist")
	}
	return c.NoContent(http.StatusOK)
}
