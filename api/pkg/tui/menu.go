package tui

import (
	"fmt"
	localnet "ragAPI/pkg/local-net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const TOTAL_OPTIONS = 4

var OPTIONS = []string{
	"Knowledge Base Menu",
	"Chat Menu",
	"Single Shot Message",
	"Quit",
}

type MainMenu struct {
	rag           *localnet.LocalControler
	focusedOption int
}

func NewMenu(rag *localnet.LocalControler) MainMenu {
	return MainMenu{
		rag,
		0,
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) ReturnFocusedOption() (tea.Model, tea.Cmd) {
	switch m.focusedOption {
	case 1:
		chat := NewChatMenu(m.rag)
		return chat, chat.Loader.Tick()
	case 2:
		return NewFirstMessageSetup(m.rag, "Single Shot Message", true), nil
	}
	return m, tea.Quit
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyDown && m.focusedOption < TOTAL_OPTIONS {
			m.focusedOption++
		}
		if msg.Type == tea.KeyUp && m.focusedOption > 0 {
			m.focusedOption--
		}
		if msg.Type == tea.KeyEnter {
			return m.ReturnFocusedOption()
		}
	}
	return m, nil
}

func (m MainMenu) View() string {
	menu := ""
	for i := range TOTAL_OPTIONS {
		line := fmt.Sprintf("%d.%s", i+1, OPTIONS[i])
		if i == m.focusedOption {
			line = lipgloss.NewStyle().Bold(true).Render(line)
		}
		menu = fmt.Sprintf("%s\n%s", menu, line)

	}
	return menu
}
