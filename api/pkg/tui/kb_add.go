package tui

import (
	"fmt"
	localnet "ragAPI/pkg/local-net"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// / Menu that displays and helps to add new files to an
// / existing kb
type KBFileAdder struct {
	rag         *localnet.LocalControler
	picker      filepicker.Model
	pickedFiles []string
	kb          string
}

func NewFileAdder(kb string, rag *localnet.LocalControler) KBFileAdder {
	return KBFileAdder{
		rag:         rag,
		kb:          kb,
		picker:      filepicker.New(),
		pickedFiles: make([]string, 0),
	}
}

func (k KBFileAdder) Init() tea.Cmd {
	return nil
}

func (k KBFileAdder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		//Return to KBMenu
		if msg.Type == tea.KeyEscape {
			return NewKBMenu(k.rag)
		}
		//Diselect a file
		if msg.Type == tea.KeyDelete && len(k.pickedFiles) > 0 {
			k.pickedFiles = k.pickedFiles[0 : len(k.pickedFiles)-1]
			return k, nil
		}
	}
	var cmd tea.Cmd
	k.picker, cmd = k.picker.Update(msg)
	if selectedFile, path := k.picker.DidSelectFile(msg); selectedFile {
		k.pickedFiles = append(k.pickedFiles, path)
	}

	return k, cmd
}

func (k KBFileAdder) pickedView() string {
	view := ""
	for index, value := range k.pickedFiles {
		view = fmt.Sprintf("%s%d.\t%s\n", view, index, value)
	}
	return view
}

func (k KBFileAdder) View() string {
	return fmt.Sprintf("Picked:\n%sSelect New:\n%s", k.pickedView(), k.picker.View())
}
