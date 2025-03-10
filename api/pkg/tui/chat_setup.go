package tui

import (
	"fmt"
	"ragAPI/pkg"
	localnet "ragAPI/pkg/local-net"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
Works as a way to set up a new chat by defining the name
*/
type ChatNameSetup struct {
	rag      *localnet.LocalControler
	chatName string
}

type ChatFirstMessageSetup struct {
	rag      *localnet.LocalControler
	header   ChatHeader
	message  string
	loading  bool
	spinType spinner.Model
	loadChat chan any
}

func NewFirstMessageSetup(
	rag *localnet.LocalControler,
	chatName string,
) ChatFirstMessageSetup {
	return ChatFirstMessageSetup{
		rag,
		ChatHeader{ChatName: chatName},
		"",
		false,
		spinner.New(
			spinner.WithSpinner(spinner.Meter),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))),
		),
		nil,
	}
}

func NewChatSetup(rag *pkg.RAG) ChatNameSetup {
	return ChatNameSetup{
		localnet.NewLocalControler(rag),
		"",
	}
}

func (c ChatNameSetup) Init() tea.Cmd {
	return nil
}

func (c ChatNameSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEscape.String():
			return c, tea.Quit
		case "enter":
			//Return a new Chat instance
			return NewFirstMessageSetup(c.rag, c.chatName), nil
		case tea.KeyBackspace.String():
			if len(c.chatName) == 0 {
				break
			}
			c.chatName = c.chatName[:len(c.chatName)-1]
		default:
			c.chatName += msg.String()
		}
	}
	return c, nil
}

func (c ChatNameSetup) View() string {
	return fmt.Sprintf("Enter chat name: %s|\n", c.chatName)
}

func (c ChatFirstMessageSetup) Init() tea.Cmd {
	return nil
}

func (c ChatFirstMessageSetup) callLoadChat() ChatFirstMessageSetup {
	c.loadChat = make(chan any)
	go func() {
		response, err := c.rag.SendNewMessageToChat(pkg.ChatInstruct{
			Message: pkg.MessageInstruct{
				Message: c.message,
			},
			NewChat:  true,
			ChatName: c.header.ChatName,
		})
		if err != nil {
			c.loadChat <- err
		} else {
			c.loadChat <- response
		}

	}()
	return c
}

func (c ChatFirstMessageSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEscape.String():
			return c, tea.Quit
		case "enter":
			//Return a new Chat instance
			c.loading = true

			return c, c.spinType.Tick
		case tea.KeyBackspace.String():
			if len(c.message) == 0 {
				break
			}
			c.message = c.message[:len(c.message)-1]
		default:
			c.message += msg.String()
		}
	case spinner.TickMsg:
		select {
		case <-c.loadChat:
			return c, tea.Quit
		default:
			newSpin, cmd := c.spinType.Update(msg)
			c.spinType = newSpin
			return c, cmd
		}
	default:
		break

	}
	return c, nil
}

func (c ChatFirstMessageSetup) View() string {
	if c.loading {
		return c.spinType.View()
	}

	return fmt.Sprintf("%s\n>> %s|\n", c.header.View(), c.message)
}
