package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
)

type (
	SubmitMsg struct {
		values map[int]string
	}
	CancelMsg struct{}
)

func cancelCmd() tea.Msg {
	return CancelMsg{}
}

type Model struct {
	inputs []textinput.Model
	title  string

	focusIndex int
	cursorMode cursor.Mode

	submitBtnIndex int
	cancelBtnIndex int
}

func New(inputs []textinput.Model, title string) Model {
	l := len(inputs)

	for i := range inputs {
		inputs[i].Cursor.Style = cursorStyle
	}

	return Model{
		inputs: inputs,
		title:  title,

		submitBtnIndex: l,
		cancelBtnIndex: l + 1,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputsLen := len(m.inputs)
	inputMaxIndex := m.cancelBtnIndex

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, cancelCmd

		case "enter":
			switch m.focusIndex {
			case m.submitBtnIndex:
				return m, func() tea.Msg {
					values := make(map[int]string)
					for i, input := range m.inputs {
						values[i] = input.Value()
					}

					return SubmitMsg{
						values: values,
					}
				}
			case m.cancelBtnIndex:
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

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < inputsLen-1 {
			b.WriteRune('\n')
		}
	}

	btnSubmit := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	if m.focusIndex == m.submitBtnIndex {
		btnSubmit = focusedStyle.Render("[ Submit ]")
	}

	btnCancel := fmt.Sprintf("[ %s ]", blurredStyle.Render("Cancel"))
	if m.focusIndex == m.cancelBtnIndex {
		btnCancel = focusedStyle.Render("[ Cancel ]")
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n\n", btnSubmit, btnCancel)

	b.WriteString(helpStyle.Render("\nPress esc to cancel"))

	return b.String()
}
