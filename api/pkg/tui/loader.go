package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoadMsg struct {
	Tck tea.Msg
	Chn chan any
}

// Base loader for the whole project
type Loader struct {
	Loading bool
	spinn   spinner.Model
	chn     chan any
	Value   any
}

func NewLoader() Loader {
	return Loader{
		false,
		spinner.New(
			spinner.WithSpinner(spinner.Meter),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))),
		),
		make(chan any),
		nil,
	}
}

func (l Loader) manageLoaderMsg(msg LoadMsg) (tea.Model, tea.Cmd) {
	select {
	case value := <-msg.Chn:
		l.Value = value
		close(l.chn)
		return l, nil
	default:
		var cmd tea.Cmd
		l.spinn, cmd = l.spinn.Update(msg.Tck)
		return l, func() tea.Msg { return LoadMsg{cmd(), l.chn} }
	}
}

func (l Loader) Tick() tea.Cmd {
	return func() tea.Msg {
		return LoadMsg{
			l.spinn.Tick(),
			l.chn,
		}
	}
}

func (l Loader) Init() tea.Cmd {
	return nil
}

func (l Loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadMsg:
		return l.manageLoaderMsg(msg)
	default:
		return l, nil
	}
}

func (l Loader) View() string {
	return l.spinn.View()
}
