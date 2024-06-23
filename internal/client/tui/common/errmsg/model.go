package errmsg

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	GotoModelMsg struct {
		Model tea.Model
	}
)

type Model struct {
	err       error
	gotoModel tea.Model
}

func New(err error, gotoModel tea.Model) *Model {
	return &Model{
		err:       err,
		gotoModel: gotoModel,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg {
			return GotoModelMsg{
				Model: m.gotoModel,
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("Error. %s\nPress any key to continue.", m.err.Error())
}
