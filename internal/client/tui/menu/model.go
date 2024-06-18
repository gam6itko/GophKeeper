package menu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
)

type Model struct {
	list list.Model
}

func NewModel(items []list.Item, ws common.WindowSize) *Model {
	m := &Model{
		list: list.New(items, list.NewDefaultDelegate(), ws.Width, ws.Height),
	}

	return m
}

func (ths Model) Init() tea.Cmd {
	return nil
}

func (ths Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ths, tea.Quit
		}

		// select option
		//if msg.Type == tea.KeyEnter {
		//	switch ths.list.SelectedItem().(common.ListItem).ID() {
		//	case loginOption:
		//		next := root.NewScreen(ths)
		//		return login.NewScreen(ths, NewLoginHandler(ths, next, &server.MockLoginHandler{})), nil
		//	case registrationOption:
		//		return login.NewScreen(ths, NewRegistrationHandler(ths, ths, &server.MockRegistrationHandler{})), nil
		//	case exitOption:
		//		return ths, tea.Quit
		//	}
		//}

	case tea.WindowSizeMsg:
		ths.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	ths.list, cmd = ths.list.Update(msg)
	return ths, cmd
}

func (ths Model) View() string {
	return ths.list.View()
}

func (ths Model) SetSize(size common.WindowSize) {
	ths.list.SetSize(size.Width, size.Height)
}
