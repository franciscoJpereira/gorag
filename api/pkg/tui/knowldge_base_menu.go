package tui

import (
	"fmt"
	"ragAPI/pkg"
	localnet "ragAPI/pkg/local-net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KnowledgeBaseMenu struct {
	rag     *localnet.LocalControler
	submenu tea.Model
}

func NewKBMenu(rag *localnet.LocalControler) (KnowledgeBaseMenu, tea.Cmd) {
	loader := NewLoader()
	kbMenu := KnowledgeBaseMenu{
		rag,
		loader,
	}
	go kbMenu.LoadChats()
	return kbMenu, loader.Tick()
}

func (k KnowledgeBaseMenu) LoadChats() {
	kbs, err := k.rag.GetAvailableKBs()
	loader, ok := k.submenu.(Loader)
	if ok && err != nil {
		loader.chn <- err
	} else if ok {
		loader.chn <- pkg.DecodeBase64Batch(kbs)
	}
}

func (k KnowledgeBaseMenu) Init() tea.Cmd {
	return nil
}

func (k KnowledgeBaseMenu) manageKeyMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	submenu, cmd := k.submenu.Update(msg)
	k.submenu = submenu
	return k, cmd
}

func (k KnowledgeBaseMenu) manageLoadMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	newSubmenu, cmd := k.submenu.Update(msg)
	loader, ok := newSubmenu.(Loader)
	if ok && loader.Value != nil {
		err, ok := loader.Value.(error)
		value, _ := loader.Value.([]string)
		if ok {
			return ErrorPopup{err.Error()}, nil
		} else {
			k.submenu = NewKBList(k.rag, value)
			return k, nil
		}
	} else {
		k.submenu = newSubmenu
		return k, cmd
	}
}

func (k KnowledgeBaseMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case LoadMsg:
		return k.manageLoadMsg(msg)
	default:
		return k.manageKeyMsg(msg)
	}
}

func (k KnowledgeBaseMenu) View() string {
	return k.submenu.View()
}

type ListKnowledgeBaseMenu struct {
	rag     *localnet.LocalControler
	bases   []string
	focused int
}

type KBMsg struct {
	Base    string
	NewBase bool
}

func NewKBCmd(newBase bool, focused string) tea.Cmd {
	return func() tea.Msg {
		return KBMsg{
			focused,
			newBase,
		}
	}
}

func NewKBList(rag *localnet.LocalControler, bases []string) ListKnowledgeBaseMenu {
	bases = append([]string{"New Base"}, bases...)
	return ListKnowledgeBaseMenu{
		rag:     rag,
		bases:   bases,
		focused: 0,
	}
}

func (l ListKnowledgeBaseMenu) Init() tea.Cmd {
	return nil
}

func (l ListKnowledgeBaseMenu) manageKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyDown:
		if l.focused < len(l.bases)-1 {
			l.focused++
		}
	case tea.KeyUp:
		if l.focused > 0 {
			l.focused--
		}
	case tea.KeyEnter:
		return l, NewKBCmd(l.focused == 0, l.bases[l.focused])
	case tea.KeyEsc:
		return NewMenu(l.rag), nil
	}
	return l, nil
}

func (l ListKnowledgeBaseMenu) manageKBMsg(msg KBMsg) (tea.Model, tea.Cmd) {
	if msg.NewBase {
		return CreateKB(l.rag), nil
	}
	return NewFileAdder(msg.Base, l.rag), nil
}

func (l ListKnowledgeBaseMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return l.manageKeyMsg(msg)
	case KBMsg:
		return l.manageKBMsg(msg)
	}

	return l, nil
}

func (l ListKnowledgeBaseMenu) View() string {
	list := ""
	for index, v := range l.bases {
		line := fmt.Sprintf("%d. %s", index+1, v)
		if index == l.focused {
			line = lipgloss.NewStyle().Bold(true).Render(line)
		}
		list += fmt.Sprintf("%s\n", line)
	}
	return list
}
