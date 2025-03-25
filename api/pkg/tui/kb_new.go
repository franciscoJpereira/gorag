package tui

import (
	"fmt"
	"ragAPI/pkg"
	localnet "ragAPI/pkg/local-net"

	tea "github.com/charmbracelet/bubbletea"
)

type NewKB struct {
	rag    *localnet.LocalControler
	loader Loader
	kbName string
}

func CreateKB(rag *localnet.LocalControler) NewKB {
	return NewKB{
		rag,
		NewLoader(),
		"",
	}
}

func (n NewKB) Init() tea.Cmd {
	return nil
}

func (n NewKB) createKB() {
	err := n.rag.CreateKB(
		pkg.EncodeBase64(n.kbName),
	)
	if err != nil {
		n.loader.chn <- err
	} else {
		n.loader.chn <- struct{}{}
	}
}

func (n NewKB) manageKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if n.kbName != "" {
			n.loader.Loading = true
			go n.createKB()
			return n, n.loader.Tick()
		}
	case tea.KeyBackspace:
		if len(n.kbName) > 0 {
			n.kbName = n.kbName[:len(n.kbName)-1]
		}
	default:
		n.kbName += msg.String()
	}
	return n, nil
}

func (n NewKB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return n.manageKeyMsg(msg)
	default:
		loader, cmd := n.loader.Update(msg)
		n.loader = loader.(Loader)
		if n.loader.Value != nil {
			return NewFileAdder(n.kbName, n.rag)
		}
		return n, cmd
	}
}

func (n NewKB) View() string {
	if n.loader.Loading {
		return n.loader.View()
	}

	return fmt.Sprintf("New KB Name\n> %s|\n", n.kbName)
}
