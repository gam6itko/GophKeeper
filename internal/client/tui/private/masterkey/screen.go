package masterkey

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Screen struct {
	prev      tea.Model
	textInput textinput.Model
	err       error
}

func NewScreen(prev tea.Model) *Screen {
	ti := textinput.New()
	ti.Placeholder = "master"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &Screen{
		prev:      prev,
		textInput: ti,
		err:       nil,
	}
}

func (ths Screen) Init() tea.Cmd {
	return textinput.Blink
}

func (ths Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return ths.prev, nil
		}

	// We handle errors just like any other message
	case errMsg:
		ths.err = msg
		return ths, nil
	}

	ths.textInput, cmd = ths.textInput.Update(msg)
	return ths, cmd
}

func (ths Screen) View() string {
	return fmt.Sprintf(
		"Enter master key?\n\n%s\n\n%s",
		ths.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
