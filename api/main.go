package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	_ "ragAPI/docs"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	knowledgebase "ragAPI/pkg/knowledge-base"
	localnet "ragAPI/pkg/local-net"
	"ragAPI/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const CONFIG_PATH = "./config/config.yaml"

func InjectRAG(r *pkg.RAG) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(pkg.RAGKey, r)
			return next(c)
		}
	}
}

func ServerMain() {
	rag := &pkg.RAG{}
	configPath, err := filepath.Abs(CONFIG_PATH)
	if err != nil {
		panic(fmt.Sprintf("Getting config path: %s\n", err))
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("Opening config file: %s\n", err))
	}
	configT, err := pkg.GetConfiguration(configFile)
	if err != nil {
		panic(fmt.Sprintf("Parsing config file: %s\n", err))
	}

	storePath, err := filepath.Abs(configT.GetStoreConfig())
	if err != nil {
		panic(fmt.Sprintf("Getting absolute path: %s", err))
	}
	e := echo.New()
	modelUrl, modelName := configT.GetModelConfig()
	rag.Api = apiinterface.NewOpenAIChatModel(
		context.Background(),
		modelUrl,
		modelName,
	)
	rag.ChatStore, err = store.NewJsonStore(storePath)
	if err != nil {
		panic(fmt.Sprintf("Creating chat store: %s", err))
	}
	rag.Kb, err = knowledgebase.NewChromaKB(
		context.Background(),
		configT.GetChromaConfig(),
	)

	e.Use(InjectRAG(rag))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/knowledge-base/:KBName", pkg.CreateKB)
	e.POST("/knowledge-base", pkg.AddDataToKB)
	e.GET("/knowledge-base", pkg.GetAvailableKBs)
	e.POST("/message", pkg.SingleShotMessage)
	e.POST("/chat", pkg.SendNewMessageToChat)
	e.GET("/chat", pkg.RetrieveAvailableChats)

	e.Logger.Fatal(e.Start(configT.GetServerConfig()))
}

func main() {
	rag := &pkg.RAG{}
	configPath, err := filepath.Abs(CONFIG_PATH)
	if err != nil {
		panic(fmt.Sprintf("Getting config path: %s\n", err))
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("Opening config file: %s\n", err))
	}
	configT, err := pkg.GetConfiguration(configFile)
	if err != nil {
		panic(fmt.Sprintf("Parsing config file: %s\n", err))
	}

	storePath, err := filepath.Abs(configT.GetStoreConfig())
	if err != nil {
		panic(fmt.Sprintf("Getting absolute path: %s", err))
	}
	modelName, modelUrl := configT.GetModelConfig()
	rag.Api = apiinterface.NewOpenAIChatModel(
		context.Background(),
		modelUrl,
		modelName,
	)
	rag.ChatStore, err = store.NewJsonStore(storePath)
	if err != nil {
		panic(fmt.Sprintf("Creating chat store: %s", err))
	}
	rag.Kb, err = knowledgebase.NewChromaKB(
		context.Background(),
		configT.GetChromaConfig(),
	)
	p := tea.NewProgram(tui.NewMenu(localnet.NewLocalControler(rag)))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error during execution: %s\n", err)
		os.Exit(1)
	}
}
