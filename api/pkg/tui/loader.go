package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Base loader for the whole project
type Loader struct {
	Loading bool
	spinn   spinner.Model
}

func NewLoader() Loader {
	return Loader{
		false,
		spinner.New(
			spinner.WithSpinner(spinner.Meter),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))),
		),
	}
}

func (l Loader) Tick() tea.Cmd {
	return l.spinn.Tick
}

func (l Loader) Init() tea.Cmd {
	return nil
}

func (l Loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	l.spinn, cmd = l.spinn.Update(msg)
	return l, cmd
}

func (l Loader) View() string {
	return l.spinn.View()
}
