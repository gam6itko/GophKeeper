package root

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/menu/common"
)

// Model управление главным меню программы.
type Model struct {
	list list.Model
}

func New(title string, width, height int) Model {
	l := list.New(
		[]list.Item{
			common.NewCmdItem("Login", "into account", gotoLoginCmd),
			common.NewCmdItem("Registration", "for a new user.", gotoRegistrationCmd),
			common.NewCmdItem("Quit", "close app", tea.Quit),
		},
		list.NewDefaultDelegate(),
		width,
		height,
	)
	l.Title = title

	return Model{
		list: l,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i := m.list.SelectedItem().(common.CmdItem)
			return m, i.Cmd()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
