package root

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
	"github.com/gam6itko/goph-keeper/internal/client/tui/menu"
)

var buildVersion = "0.0.0"
var buildDate = "never"

const (
	loginOption = iota
	registrationOption
	exitOption
)

var windowSize common.WindowSize

type gotoMainMenuMsg struct{}
type exitMsg struct{}

type Model struct {
	initialized bool
	stack       []tea.Model
	index       int
}

func NewModel() *Model {
	//items :=

	//m := &Model{
	//	list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	//}
	//m.list.Title = fmt.Sprintf("GophKeeper. Version: %s. Build: %s.", buildVersion, buildDate)

	return &Model{
		stack: make([]tea.Model, 0, 4),
		index: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowSize = common.WindowSize{Width: msg.Width, Height: msg.Height}
		if m.initialized {
			//to
		} else {
			return m, func() tea.Msg { return gotoMainMenuMsg{} }
		}
	case gotoMainMenuMsg:
		m.index = 0
		m.stack = []tea.Model{
			menu.NewModel(
				[]list.Item{
					menu.NewItem("Login", "into account", func() tea.Msg {
						return nil
					}),
					menu.NewItem("Registration", "for a new user.", func() tea.Msg {
						return nil
					}),
					menu.NewItem("Exit", "", func() tea.Msg {
						return exitMsg{}
					}),
				},
				windowSize,
			),
		}
	case exitMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// select option
		//if msg.Type == tea.KeyEnter {
		//	switch m.list.SelectedItem().(common.ListItem).ID() {
		//	case loginOption:
		//		next := root.NewScreen(m)
		//		return login.NewScreen(m, NewLoginHandler(m, next, &server.MockLoginHandler{})), nil
		//	case registrationOption:
		//		return login.NewScreen(m, NewRegistrationHandler(m, m, &server.MockRegistrationHandler{})), nil
		//	case exitOption:
		//		return m, tea.Quit
		//	}
		//}

	}

	if model := m.topModel(); model == nil {
		return m, nil
	}

	return m.topModel().Update(msg)
}

func (m Model) View() string {
	model := m.topModel()
	if model == nil {
		return ""
	}
	return model.View()
}

func (m Model) topModel() tea.Model {
	if m.index == -1 {
		return nil
	}

	return m.stack[m.index]
}
