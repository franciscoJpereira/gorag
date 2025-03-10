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

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return func() tea.Msg {
		return "Hola, Init!"
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(tui.NewChatSetup(nil))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error during execution: %s\n", err)
		os.Exit(1)
	}
}
