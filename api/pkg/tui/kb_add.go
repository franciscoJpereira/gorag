package tui

import (
	"fmt"
	"os"
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
	lastEnter   bool
}

func NewFileAdder(kb string, rag *localnet.LocalControler) (tea.Model, tea.Cmd) {
	var err error
	picker := filepicker.New()
	picker.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	picker.Height = 10
	kbFilePicker := KBFileAdder{
		rag:         rag,
		kb:          kb,
		picker:      picker,
		pickedFiles: make([]string, 0),
	}
	return kbFilePicker, kbFilePicker.Init()
}

func (k KBFileAdder) Init() tea.Cmd {
	return k.picker.Init()
}

func (k KBFileAdder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		//Return to KBMenu
		if msg.Type == tea.KeyEscape {
			return NewKBMenu(k.rag)
		}
		//Diselect a file
		if msg.Type == tea.KeyBackspace && len(k.pickedFiles) > 0 {
			k.pickedFiles = k.pickedFiles[0 : len(k.pickedFiles)-1]
			return k, nil
		}

		if msg.Type == tea.KeyEnter {
			if k.lastEnter {
				//TODO: Read the files and send the data to the back end
				return NewMenu(k.rag), nil
			}
			k.lastEnter = true
		} else {
			k.lastEnter = false
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
		view = fmt.Sprintf("%s%d. %s\n", view, index+1, value)
	}
	return view
}

func (k KBFileAdder) View() string {
	return fmt.Sprintf("Picked:\n%s\nSelect New:\n%s", k.pickedView(), k.picker.View())
}
