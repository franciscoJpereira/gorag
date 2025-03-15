package tui

import (
	"fmt"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	localnet "ragAPI/pkg/local-net"

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
		singleShot,
	}
}

func NewChatSetup(rag *localnet.LocalControler) ChatNameSetup {
	return ChatNameSetup{
		rag,
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
		c.loader.chn <- err
	} else {
		c.loader.chn <- response
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
		c.loader.chn <- err
	} else {
		c.loader.chn <- response
	}
}

func (c ChatFirstMessageSetup) callLoadChat() ChatFirstMessageSetup {
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

func (c ChatFirstMessageSetup) handleTickMessage(msg LoadMsg) (tea.Model, tea.Cmd) {
	loader, cmd := c.loader.Update(msg)
	c.loader = loader.(Loader)
	if c.loader.Value == nil {
		return c, cmd
	}
	err, ok := c.loader.Value.(error)
	value, _ := c.loader.Value.(pkg.MessageResponse)
	if ok {
		return ErrorPopup{err.Error()}, nil
	} else {
		newChat := NewChat(
			c.rag,
			store.ChatHistory{
				ChatName: c.header.ChatName,
				Messages: []apiinterface.ChatMessage{
					{
						Role:    "user",
						Content: value.Query,
					},
					{
						Role:    "system",
						Content: value.Response,
					},
				},
			},
		)
		return newChat, nil
	}
}

func (c ChatFirstMessageSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyMsg(msg)
	case LoadMsg:
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
