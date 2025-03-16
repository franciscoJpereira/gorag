package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DEFAULT_WIDTH           = 80
	DEFAULT_HEIGHT          = 10
	DEFAULT_VERTICAL_SCROLL = 5
)

type View struct {
	originalText string
	textLines    []string
	firstShown   int
	width        int
	height       int
	totalDowns   int
}

func (c View) makeLine(initial int) (line string, last int) {
	actual := initial
	for actual < len(c.originalText) && actual-initial < c.width {
		if c.originalText[actual] == '\n' {
			break
		}
		line += string(c.originalText[actual])
		actual++
	}
	last = actual + 1
	return
}

func (c View) makeLines() View {
	initial := 0
	for initial < len(c.originalText) {
		line, last := c.makeLine(initial)
		initial = last
		c.textLines = append(c.textLines, line)
	}
	return c
}

func NewView(text string) View {
	view := View{
		originalText: text,
		textLines:    []string{},
		firstShown:   0,
		width:        DEFAULT_WIDTH,
		height:       DEFAULT_HEIGHT,
		totalDowns:   0,
	}
	return view.makeLines()
}

func (c View) Init() tea.Cmd {
	return nil
}

func (c View) manageKeyUp() (tea.Model, tea.Cmd) {
	c.firstShown = max(0, c.firstShown-DEFAULT_VERTICAL_SCROLL)
	return c, nil
}

func (c View) manageKeyDown() (tea.Model, tea.Cmd) {
	c.firstShown = min(len(c.textLines)-c.height, c.firstShown+DEFAULT_VERTICAL_SCROLL)
	return c, nil
}

func (c View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			return c.manageKeyUp()
		case tea.KeyDown:
			return c.manageKeyDown()
		}
	}
	return c, nil
}

func (c View) makeView() string {
	if c.firstShown+c.height >= len(c.textLines) {
		return strings.Join(c.textLines[c.firstShown:], "\n")
	}

	return strings.Join(c.textLines[c.firstShown:c.firstShown+c.height], "\n")
}

func (c View) View() string {
	return c.makeView()
}
