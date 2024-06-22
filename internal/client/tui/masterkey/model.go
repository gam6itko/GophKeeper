package masterkey

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
)

type (
	SubmitMsg struct {
		Key      []byte
		RetryMsg tea.Msg
	}
	CancelMsg struct {
		Prev tea.Model
	}
)

// Model form decorator.
type Model struct {
	form tea.Model

	successRetryMsg tea.Msg
	prev            tea.Model
}

func New(form tea.Model, successRetryMsg tea.Msg, prev tea.Model) *Model {
	return &Model{
		form: form,

		successRetryMsg: successRetryMsg,
		prev:            prev,
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case form.SubmitMsg:
		return m, func() tea.Msg {
			return SubmitMsg{
				Key:      []byte(msg.Values[0]),
				RetryMsg: m.successRetryMsg,
			}
		}
	case form.CancelMsg:
		return m, func() tea.Msg {
			return CancelMsg{Prev: m.prev}
		}
	}

	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.form.View()
}
