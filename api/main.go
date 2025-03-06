package main

import (
	"context"
	"fmt"
	"path/filepath"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	knowledgebase "ragAPI/pkg/knowledge-base"

	_ "ragAPI/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func TestMessage(message string) {
	service := apiinterface.NewOpenAIChatModel(
		context.Background(),
		"http://localhost:1234/v1/",
		"deepseek-r1-distill-llama-8b@q6_k")
	response, err := service.Send(message)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Response: %s\n", response.Content)
	}
	models, _ := service.Models()
	fmt.Printf("Models: %v\n", models)
}

func TestChroma() {
	chromaOptions := knowledgebase.ChromaKBOptions{
		BasePath:      "http://localhost:8000",
		EmbedderPath:  "http://localhost:1234/v1/embeddings",
		EmbedderModel: "text-embedding-granite-embedding-278m-multilingual",
		MaxResults:    10,
	}
	chromadb, err := knowledgebase.NewChromaKB(context.Background(), chromaOptions)
	if err != nil {
		panic(err)
	}
	if err = chromadb.CreateColletion("TestingCollection"); err != nil {
		panic(fmt.Sprintf("Creating collection: %s", err))
	}
	if err = chromadb.AddDataToCollection("TestingCollection", []string{"South park es muy divertido", "Un cuento de sueños locos", "Un perro caminando por la calle no es bueno"}); err != nil {
		panic(fmt.Sprintf("Adding data to collection: %s", err))
	}
	result := chromadb.Retrieve("TestingCollection", "Perros")
	fmt.Println(result)
}

func SimpleRPL() {
	chromaOptions := knowledgebase.ChromaKBOptions{
		BasePath:      "http://localhost:8000",
		EmbedderPath:  "http://localhost:1234/v1/embeddings",
		EmbedderModel: "text-embedding-granite-embedding-278m-multilingual",
		MaxResults:    10,
	}
	chromadb, err := knowledgebase.NewChromaKB(context.Background(), chromaOptions)
	if err != nil {
		panic(fmt.Sprintf("Creating chromadb: %s", err))
	}
	apiClient := apiinterface.NewOpenAIChatModel(
		context.Background(),
		"http://localhost:1234/v1/",
		"deepseek-r1-distill-llama-8b@q6_k")
	path, err := filepath.Abs("../store/")
	if err != nil {
		panic(fmt.Sprintf("Getting absolute path: %s", err))
	}
	chatStore, err := store.NewJsonStore(path)
	if err != nil {
		panic(fmt.Sprintf("Creating chat store: %s", err))
	}
	rag := pkg.RAG{
		Kb:        chromadb,
		Api:       apiClient,
		ChatStore: chatStore,
	}

	for {
		var message string
		fmt.Scanln(&message)
		r, err := rag.SingleShotMessage(pkg.MessageInstruct{
			Message: message,
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("Response: %s\n", r)
		}
	}
}

const RAGKey = "RAG"

func InjectRAG(r *pkg.RAG) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(RAGKey, r)
			return next(c)
		}
	}
}

func main() {
	//TestMessage("Hello there!")
	//TestChroma()
	//SimpleRPL()
	rag := &pkg.RAG{}
	storePath, err := filepath.Abs("../store/")
	if err != nil {
		panic(fmt.Sprintf("Getting absolute path: %s", err))
	}
	e := echo.New()
	rag.Api = apiinterface.NewOpenAIChatModel(
		context.Background(),
		"http://localhost:1234/v1/",
		"deepseek-r1-distill-llama-8b@q6_k",
	)
	rag.ChatStore, err = store.NewJsonStore(storePath)
	if err != nil {
		panic(fmt.Sprintf("Creating chat store: %s", err))
	}
	rag.Kb, err = knowledgebase.NewChromaKB(
		context.Background(),
		knowledgebase.ChromaKBOptions{
			BasePath:      "http://localhost:8000",
			EmbedderPath:  "http://localhost:1234/v1/embeddings",
			EmbedderModel: "text-embedding-granite-embedding-278m-multilingual",
			MaxResults:    10,
		},
	)

	e.Use(InjectRAG(rag))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/knowledge-base/:KBName", CreateKB)
	e.POST("/knowledge-base", AddDataToKB)
	e.GET("/knowledge-base", GetAvailableKBs)
	e.POST("/message", SingleShotMessage)
	e.POST("/chat", SendNewMessageToChat)
	e.GET("/chat", RetrieveAvailableChats)
	e.Logger.Fatal(e.Start(":1323"))
}
