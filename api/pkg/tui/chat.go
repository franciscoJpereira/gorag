package tui

import (
	"fmt"
	"ragAPI/pkg"
	apiinterface "ragAPI/pkg/apiInterface"
	"ragAPI/pkg/chat/store"
	localnet "ragAPI/pkg/local-net"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// Shows the user question, query, prompt and the corresponding model response
type ChatPiece struct {
	header    ChatHeader
	UserQuery string
	Response  string
	port      viewport.Model
}

func generateContent(query, response string, port viewport.Model) string {
	user := fmt.Sprintf("User >> \n%s\n", query)
	model := fmt.Sprintf("Model >> \n%s\n", response)
	line := strings.Repeat("=", max(0, port.Width))
	return lipgloss.JoinVertical(lipgloss.Left, user, line, model)
}

func NewChatPiece(header ChatHeader, query string, response string) ChatPiece {
	port := viewport.New(100, 10)
	port.SetContent(generateContent(query, response, port))
	return ChatPiece{
		header,
		query,
		response,
		port,
	}
}

func (c ChatPiece) Init() tea.Cmd {
	return nil
}

func (c ChatPiece) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return c, tea.Quit
		}
	case tea.WindowSizeMsg:
		c.port.Height = msg.Height
		c.port.Width = msg.Width
	}
	c.port, cmd = c.port.Update(msg)

	return c, tea.Batch([]tea.Cmd{cmd}...)
}

func (c ChatPiece) View() string {
	line := strings.Repeat("-", max(0, c.port.Width))
	trail := strings.Repeat("=", max(0, c.port.Width))
	return lipgloss.JoinVertical(lipgloss.Left, c.header.View(), line, c.port.View(), trail)
}

type ChatMessage struct {
	mesage string
}

func (c ChatMessage) Init() tea.Cmd {
	return nil
}

func (c ChatMessage) manageKeyMsg(msg tea.KeyMsg) ChatMessage {
	switch msg.Type {
	case tea.KeyBackspace:
		if len(c.mesage) > 0 {
			c.mesage = c.mesage[:len(c.mesage)-1]
		}
	default:
		c.mesage += msg.String()
	}
	return c
}

func (c ChatMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		c = c.manageKeyMsg(msg)
	}
	return c, nil
}

func (c ChatMessage) View() string {
	return fmt.Sprintf(">> %s|\n", c.mesage)
}

type Chat struct {
	rag          *localnet.LocalControler
	history      store.ChatHistory
	current      ChatPiece
	currentIndex int
	message      ChatMessage
	loader       Loader
}

func NewChat(rag *localnet.LocalControler, history store.ChatHistory) Chat {
	currentIndex := len(history.Messages) - 2
	firstPiece := NewChatPiece(ChatHeader{history.ChatName}, history.Messages[currentIndex].Content, history.Messages[currentIndex+1].Content)
	return Chat{
		rag:          rag,
		history:      history,
		current:      firstPiece,
		currentIndex: currentIndex,
		message:      ChatMessage{},
		loader:       NewLoader(),
	}
}

func (c Chat) Init() tea.Cmd {
	return nil
}

func (c Chat) loadMessage() {
	response, err := c.rag.SendNewMessageToChat(pkg.ChatInstruct{
		Message: pkg.MessageInstruct{
			Message: c.message.mesage,
		},
		ChatName: c.history.ChatName,
	})
	if err != nil {
		c.loader.chn <- err
	} else {

		c.loader.chn <- response
	}
}

func (c Chat) manageKeyMsg(msg tea.KeyMsg) Chat {
	switch msg.Type {
	case tea.KeyLeft:
		if c.currentIndex > 1 {
			c.currentIndex -= 2
		}
	case tea.KeyRight:
		if c.currentIndex < len(c.history.Messages)-4 {
			c.currentIndex += 2
		}
	case tea.KeyDown:
		break
	case tea.KeyUp:
		break
	default:
		message, _ := c.message.Update(msg)
		c.message = message.(ChatMessage)
	}
	c.current = NewChatPiece(
		ChatHeader{c.history.ChatName},
		c.history.Messages[c.currentIndex].Content,
		c.history.Messages[c.currentIndex+1].Content,
	)
	return c
}

func (c Chat) manageLoadMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	loader, cmd := c.loader.Update(msg)
	c.loader = loader.(Loader)
	if c.loader.Value != nil {
		c.message = ChatMessage{}
		err, ok := c.loader.Value.(error)
		response, _ := c.loader.Value.(pkg.MessageResponse)
		c.loader.Value = nil
		c.loader.Loading = false
		if ok {
			return ErrorPopup{err.Error()}, nil
		}
		c.history.Messages = append(c.history.Messages,
			apiinterface.ChatMessage{
				Role:    "user",
				Content: c.message.mesage,
			},
			apiinterface.ChatMessage{
				Role:    "system",
				Content: response.Response,
			},
		)
		c.currentIndex += 2
	}
	return c, cmd
}

func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && c.message.mesage != "" {
			c.loader.Loading = true
			go c.loadMessage()
			return c, c.loader.Tick()
		}
		c = c.manageKeyMsg(msg)
		if msg.Type == tea.KeyEscape {
			return NewMenu(c.rag), nil
		}
	default:
		return c.manageLoadMsg(msg)
	}
	current, _ := c.current.Update(msg)
	c.current = current.(ChatPiece)
	return c, nil
}

func (c Chat) View() string {
	view := c.current.View() + "\n"
	if !c.loader.Loading {
		view += c.message.View()
	} else {
		view += c.loader.View()
	}
	return view
}
