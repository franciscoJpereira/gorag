package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
