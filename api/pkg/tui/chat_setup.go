package tui

import (
	"fmt"
	"ragAPI/pkg"
	localnet "ragAPI/pkg/local-net"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

/*
Works as a way to set up a new chat by defining the name
*/
type ChatNameSetup struct {
	rag      *localnet.LocalControler
	chatName string
}

type ChatFirstMessageSetup struct {
	rag          *localnet.LocalControler
	header       ChatHeader
	message      string
	loader       Loader
	loadChat     chan any
	isSingleShot bool
}

func NewFirstMessageSetup(
	rag *localnet.LocalControler,
	chatName string,
	singleShot bool,
) ChatFirstMessageSetup {

	return ChatFirstMessageSetup{
		rag,
		ChatHeader{ChatName: chatName},
		"",
		NewLoader(),
		nil,
		singleShot,
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
			return NewFirstMessageSetup(c.rag, c.chatName, false), nil
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

func (c ChatFirstMessageSetup) SingleShotMessage() {
	response, err := c.rag.SingleShotMessage(pkg.MessageInstruct{
		Message: c.message,
	},
	)
	if err != nil {
		c.loadChat <- err
	} else {
		c.loadChat <- response
	}
}

func (c ChatFirstMessageSetup) NewChatMessage() {
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
}

func (c ChatFirstMessageSetup) callLoadChat() ChatFirstMessageSetup {
	c.loadChat = make(chan any)
	if c.isSingleShot {
		go c.SingleShotMessage()
	} else {
		go c.NewChatMessage()
	}
	return c
}

func (c ChatFirstMessageSetup) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case tea.KeyEscape.String():
		return c, tea.Quit
	case "enter":
		//Return a new Chat instance
		c.loader.Loading = true
		c = c.callLoadChat()
		return c, c.loader.Tick()
	case tea.KeyBackspace.String():
		if len(c.message) == 0 {
			break
		}
		c.message = c.message[:len(c.message)-1]
	default:
		c.message += msg.String()
	}
	return c, nil
}

func (c ChatFirstMessageSetup) handleTickMessage(msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	select {
	case response := <-c.loadChat:
		err, ok := response.(error)
		value, _ := response.(pkg.MessageResponse)
		if ok {
			return ErrorPopup{err.Error()}, nil
		} else {
			return NewChatPiece(
				c.header,
				value.Query,
				value.Response,
			), nil
		}
	default:
		newSpin, cmd := c.loader.Update(msg)
		c.loader = newSpin.(Loader)
		return c, cmd
	}

}

func (c ChatFirstMessageSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyMsg(msg)
	case spinner.TickMsg:
		return c.handleTickMessage(msg)
	default:
		break
	}
	return c, nil
}

func (c ChatFirstMessageSetup) View() string {
	if c.loader.Loading {
		return c.loader.View()
	}

	return fmt.Sprintf("%s\n>> %s|\n", c.header.View(), c.message)
}
