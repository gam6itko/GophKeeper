package root

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
)

var buildVersion = "0.0.0"
var buildDate = "never"

const (
	loginOption = iota
	registrationOption
	exitOption
)

var windowSize common.WindowSize

//type exitMsg struct{}

type appState int

const (
	stateIdle appState = iota
	stateOnRootMenu
	stateOnLoginScreen
	stateOnRegistrationScreen
	stateStorePrivateData
	stateLoadPrivate
)

type Model struct {
	width, height int

	current tea.Model
	state   appState
}

func New() *Model {
	return &Model{
		state: stateIdle,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.current == nil {
			if child, ok := m.current.(common.IWindowSizeAware); ok {
				child.SetSize(m.width, m.height)
			}
		}
		if m.state == stateIdle {
			return m, gotoRootMenuCmd
		}

	case gotoRootMenuMsg:
		m.state = stateOnRootMenu
		m.current = newRootMenu(
			fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate),
			m.width,
			m.height,
		)

	case gotoLoginMsg:
		m.state = stateOnLoginScreen
		m.current = newLoginForm()
		return m, m.current.Init()

	case gotoRegistrationMsg:
		m.state = stateOnRegistrationScreen
		m.current = newRegistrationForm()
		return m, m.current.Init()

	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.current == nil {
		return "Loading..."
	}

	return m.current.View()
}

//
//type Model struct {
//	initialized bool
//	сгккуте     []tea.Model
//	index       int
//}
//
//func NewModel() *Model {
//	return &Model{
//		stack: make([]tea.Model, 0, 4),
//		index: -1,
//	}
//}
//
//func (m Model) Init() tea.Cmd {
//	return nil
//}
//
//func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.WindowSizeMsg:
//		windowSize = common.WindowSize{Width: msg.Width, Height: msg.Height}
//		if m.initialized {
//			//to
//		} else {
//			return m, func() tea.Msg { return gotoRootMenuMsg{} }
//		}
//	case gotoRootMenuMsg:
//		m.index = 0
//		m.stack = []tea.Model{
//			menu.NewModel(
//				[]list.common{
//					menu.newItem("Login", "into account", func() tea.Msg {
//						return nil
//					}),
//					menu.newItem("Registration", "for a new user.", func() tea.Msg {
//						return nil
//					}),
//					menu.newItem("Exit", "", func() tea.Msg {
//						return exitMsg{}
//					}),
//				},
//				windowSize,
//			),
//		}
//	case exitMsg:
//		return m, tea.Quit
//
//	case tea.KeyMsg:
//		if msg.String() == "ctrl+c" {
//			return m, tea.Quit
//		}
//
//		// select option
//		//if msg.Type == tea.KeyEnter {
//		//	switch m.list.SelectedItem().(common.ListItem).ID() {
//		//	case loginOption:
//		//		next := root.NewScreen(m)
//		//		return login.NewScreen(m, NewLoginHandler(m, next, &server.MockLoginHandler{})), nil
//		//	case registrationOption:
//		//		return login.NewScreen(m, NewRegistrationHandler(m, m, &server.MockRegistrationHandler{})), nil
//		//	case exitOption:
//		//		return m, tea.Quit
//		//	}
//		//}
//
//	}
//
//	if model := m.topModel(); model == nil {
//		return m, nil
//	}
//
//	return m.topModel().Update(msg)
//}
//
//func (m Model) View() string {
//	model := m.topModel()
//	if model == nil {
//		return ""
//	}
//	return model.View()
//}
