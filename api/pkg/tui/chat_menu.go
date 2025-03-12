package tui

import (
	"fmt"
	localnet "ragAPI/pkg/local-net"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatMenu struct {
	rag         *localnet.LocalControler
	chats       []string
	Loader      Loader
	focusedChat int
	chn         chan any
}

func (c ChatMenu) LoadChats() {
	c.Loader.Loading = true
	chats, err := c.rag.RetrieveAvailableChats()
	if err != nil {
		c.chn <- err
	} else {
		c.chn <- chats
	}
}

func NewChatMenu(rag *localnet.LocalControler) ChatMenu {
	chat := ChatMenu{
		rag,
		[]string{"New Chat"},
		NewLoader(),
		0,
		make(chan any),
	}
	chat.Loader.Loading = true
	go chat.LoadChats()
	return chat
}

func (c ChatMenu) Init() tea.Cmd {
	return nil
}

func (c ChatMenu) manageKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if c.focusedChat > 0 {
			c.focusedChat--
		}
	case tea.KeyDown:
		if c.focusedChat < len(c.chats) {
			c.focusedChat++
		}
	case tea.KeyEnter:
		if c.focusedChat == 0 {
			return NewChatSetup(c.rag), nil
		}
		return c, tea.Quit
	case tea.KeyEsc:
		return NewMenu(c.rag), nil
	}
	return c, nil
}

func (c ChatMenu) manageMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !c.Loader.Loading {
			return c.manageKeyMsg(msg)
		}
	case spinner.TickMsg:
		loader, cmd := c.Loader.Update(msg)
		c.Loader = loader.(Loader)
		return c, cmd
	}
	return c, nil
}

func (c ChatMenu) manageLoadValue(value any) (tea.Model, tea.Cmd) {
	err, ok := value.(error)
	chats, _ := value.([]string)
	if ok {
		return ErrorPopup{err.Error()}, nil
	} else {
		c.chats = append(c.chats, chats...)
		c.Loader.Loading = false
		return c, nil
	}
}

func (c ChatMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	select {
	case value := <-c.chn:
		return c.manageLoadValue(value)
	default:
		return c.manageMsg(msg)
	}

}

func (c ChatMenu) View() string {
	if c.Loader.Loading {
		return c.Loader.View()
	}
	menu := ""
	for index, value := range c.chats {
		line := fmt.Sprintf("%d. %s", index+1, value)
		if index == c.focusedChat {
			line = lipgloss.NewStyle().Bold(true).Render(line)
		}
		menu = fmt.Sprintf("%s\n%s", menu, line)
	}
	return menu
}
