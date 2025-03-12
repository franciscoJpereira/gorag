package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Header that shows the name of the chat
type ChatHeader struct {
	ChatName string
}

func (c ChatHeader) Init() tea.Cmd {
	return nil
}

func (c ChatHeader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c ChatHeader) View() string {
	return fmt.Sprintf("Chat: %s\n", c.ChatName)
}

// Shows the user question, query, prompt and the corresponding model response
type ChatPiece struct {
	header    ChatHeader
	UserQuery string
	Response  string
	port      viewport.Model
}

func generateContent(query, response string, port viewport.Model) string {
	user := fmt.Sprintf("User >> \n%s\n", query)
	model := fmt.Sprintf("Model >> \n%s\n", response)
	line := strings.Repeat("=", max(0, port.Width))
	return lipgloss.JoinVertical(lipgloss.Left, user, line, model)
}

func NewChatPiece(header ChatHeader, query string, response string) ChatPiece {
	port := viewport.New(100, 10)
	port.SetContent(generateContent(query, response, port))
	return ChatPiece{
		header,
		query,
		response,
		port,
	}
}

func (c ChatPiece) Init() tea.Cmd {
	return nil
}

func (c ChatPiece) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return c, tea.Quit
		}
	case tea.WindowSizeMsg:
		c.port.Height = msg.Height
		c.port.Width = msg.Width
	}
	c.port, cmd = c.port.Update(msg)

	return c, tea.Batch([]tea.Cmd{cmd}...)
}

func (c ChatPiece) View() string {
	line := strings.Repeat("-", max(0, c.port.Width))
	return lipgloss.JoinVertical(lipgloss.Left, c.header.View(), line, c.port.View())
}
