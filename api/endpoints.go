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
