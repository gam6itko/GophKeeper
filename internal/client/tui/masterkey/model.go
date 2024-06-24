package masterkey

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
)

type (
	SubmitMsg struct {
		Key      []byte
		RetryCmd tea.Cmd
	}
	CancelMsg struct {
		Prev tea.Model
	}
)

// Model form decorator.
type Model struct {
	form tea.Model
	// retryCmd сообщение которое нужно послать еще раз при установке мастер ключа.
	retryCmd tea.Cmd
	prev     tea.Model
}

func New(form tea.Model, successRetryCmd tea.Cmd, prev tea.Model) *Model {
	return &Model{
		form: form,

		retryCmd: successRetryCmd,
		prev:     prev,
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
				RetryCmd: m.retryCmd,
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
