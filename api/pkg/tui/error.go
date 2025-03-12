package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// / Shows an error until the user presses enter
// / It returns to the main menu
type ErrorPopup struct {
	Error string
}

func (e ErrorPopup) wrapError(style lipgloss.Style) string {
	return fmt.Sprintf("%sPress Enter to continue",
		style.Render(
			fmt.Sprintf("Error: %s\n", e.Error),
		))
}

func (e ErrorPopup) applyStyle() string {
	return e.wrapError(
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff0000")),
	)
}

func (e ErrorPopup) Init() tea.Cmd {
	return nil
}

func (e ErrorPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			//TODO: DO not quit, go back to main menu
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e ErrorPopup) View() string {
	return e.applyStyle()
}
