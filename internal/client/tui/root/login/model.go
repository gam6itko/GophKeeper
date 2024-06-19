package login

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	titleStyle          = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

const (
	usernameInputIndex = iota
	passwordInputIndex
	submitButtonIndex
	cancelButtonIndex
)
const inputMaxIndex = cancelButtonIndex

type (
	SubmitMsg struct {
		values map[string]string
	}
	CancelMsg struct{}
)

func cancelCmd() tea.Msg {
	return CancelMsg{}
}

type Model struct {
	inputs     []textinput.Model
	focusIndex int
	cursorMode cursor.Mode

	Title string
}

func New() Model {
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

	return Model{
		inputs: []textinput.Model{
			username,
			password,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputsLen := len(m.inputs)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, cancelCmd

		case "enter":
			switch m.focusIndex {
			case submitButtonIndex:
				return m, func() tea.Msg {
					values := make(map[string]string)
					for _, input := range m.inputs {
						values[input.Na]
					}
					return SubmitMsg{
						login:    m.inputs[usernameInputIndex].Value(),
						password: m.inputs[passwordInputIndex].Value(),
					}
				}
			case cancelButtonIndex:
				return m, cancelCmd
			}

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > inputMaxIndex {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = inputsLen
			}

			cmds := make([]tea.Cmd, inputsLen)
			for i := 0; i < inputsLen; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	return m, nil
}

func (m Model) View() string {
	inputsLen := len(m.inputs)

	var b strings.Builder

	b.WriteString(titleStyle.Render(m.Title))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < inputsLen-1 {
			b.WriteRune('\n')
		}
	}

	btnSubmit := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	if m.focusIndex == submitButtonIndex {
		btnSubmit = focusedStyle.Render("[ Submit ]")
	}

	btnCancel := fmt.Sprintf("[ %s ]", blurredStyle.Render("Cancel"))
	if m.focusIndex == cancelButtonIndex {
		btnCancel = focusedStyle.Render("[ Cancel ]")
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", btnSubmit, btnCancel)

	b.WriteString(helpStyle.Render("\nPress esc to cancel"))

	return b.String()
}
