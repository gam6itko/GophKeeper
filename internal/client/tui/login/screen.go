package login

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const (
	UsernameInputIndex = iota
	PasswordInputIndex
	SubmitButtonIndex
	CancelButtonIndex
)
const InputMaxIndex = CancelButtonIndex

var (
	titleStyle          = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type Screen struct {
	prev          tea.Model
	submitHandler ISubmitHandler

	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode

	Title string
}

func NewScreen(prev tea.Model, submitHandler ISubmitHandler) *Screen {
	username := textinput.New()
	username.Placeholder = "Username"
	username.Cursor.Style = cursorStyle
	username.CharLimit = 32
	username.CharLimit = 64
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Password"
	password.Cursor.Style = cursorStyle
	password.CharLimit = 32
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	m := &Screen{
		prev:          prev,
		submitHandler: submitHandler,

		inputs: []textinput.Model{
			username,
			password,
		},
	}

	return m
}

func (ths *Screen) Init() tea.Cmd {
	return textinput.Blink
}

func (ths *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputsLen := len(ths.inputs)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return ths, tea.Quit

		case "esc":
			return ths.prev, nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				switch ths.focusIndex {
				case SubmitButtonIndex:
					return ths.submitHandler.Handle(
						ths.inputs[UsernameInputIndex].Value(),
						ths.inputs[PasswordInputIndex].Value(),
					)
				case CancelButtonIndex:
					return ths.prev, nil
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				ths.focusIndex--
			} else {
				ths.focusIndex++
			}

			if ths.focusIndex > InputMaxIndex {
				ths.focusIndex = 0
			} else if ths.focusIndex < 0 {
				ths.focusIndex = inputsLen
			}

			cmds := make([]tea.Cmd, inputsLen)
			for i := 0; i < inputsLen; i++ {
				if i == ths.focusIndex {
					// Set focused state
					cmds[i] = ths.inputs[i].Focus()
					ths.inputs[i].PromptStyle = focusedStyle
					ths.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				ths.inputs[i].Blur()
				ths.inputs[i].PromptStyle = noStyle
				ths.inputs[i].TextStyle = noStyle
			}

			return ths, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := ths.updateInputs(msg)

	return ths, cmd
}

func (ths *Screen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(ths.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range ths.inputs {
		ths.inputs[i], cmds[i] = ths.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (ths *Screen) View() string {
	inputsLen := len(ths.inputs)

	var b strings.Builder

	b.WriteString(titleStyle.Render(ths.Title))
	b.WriteString("\n\n")

	for i := range ths.inputs {
		b.WriteString(ths.inputs[i].View())
		if i < inputsLen-1 {
			b.WriteRune('\n')
		}
	}

	btnSubmit := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	if ths.focusIndex == SubmitButtonIndex {
		btnSubmit = focusedStyle.Render("[ Submit ]")
	}

	btnCancel := fmt.Sprintf("[ %s ]", blurredStyle.Render("Cancel"))
	if ths.focusIndex == CancelButtonIndex {
		btnCancel = focusedStyle.Render("[ Cancel ]")
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", btnSubmit, btnCancel)

	b.WriteString(helpStyle.Render("\nPress esc to cancel"))

	return b.String()
}
