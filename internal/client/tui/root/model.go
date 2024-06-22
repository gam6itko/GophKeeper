package root

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/loading"
	"log"
)

var buildVersion = "0.0.0"
var buildDate = "never"

const (
	loginOption = iota
	registrationOption
	exitOption
)

var windowSize common.WindowSize

type appState int

const (
	stateIdle appState = iota
	stateOnRootMenu
	stateOnLoginFrom
	stateOnRegistrationForm
	stateStorePrivateData
	stateLoadPrivate
)

type Model struct {
	width, height int

	prev    tea.Model
	current tea.Model
	state   appState

	server     server.IServer
	cancelFunc *context.CancelFunc
}

func New(server server.IServer) *Model {
	return &Model{
		state:   stateIdle,
		current: newRootMenu(fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate), 0, 0),

		server: server,
	}
}

func (m Model) Init() tea.Cmd {
	return gotoRootMenuCmd
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

	case gotoRootMenuMsg:
		m.state = stateOnRootMenu
		m.current = newRootMenu(
			fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate),
			m.width,
			m.height,
		)
		return m, m.current.Init()

	case gotoLoginMsg:
		m.state = stateOnLoginFrom
		m.current = newLoginForm()
		return m, m.current.Init()

	case gotoRegistrationMsg:
		m.state = stateOnRegistrationForm
		m.current = newRegistrationForm()
		return m, m.current.Init()

	case form.SubmitMsg:
		ctx, fnCancel := context.WithCancel(context.Background())
		m.cancelFunc = &fnCancel

		switch m.state {
		case stateOnLoginFrom:
			model := loading.New(
				func() tea.Msg {
					err := m.server.Login(
						ctx,
						server.LoginDTO{
							Username: msg.Values[LoginFormUsernameIndex],
							Password: msg.Values[LoginFormPasswordIndex],
						},
					)
					return loading.LoginResponseMsg{Err: err}
				},
				m.current,
			)
			m.prev = m.current
			m.current = model
			return m, m.current.Init()

			log.Printf("submit login")
			//todo check user-pass
			//todo server.sendLogin
			//todo hadle err or  goto private
		case stateOnRegistrationForm:
			log.Printf("submit reg")
			//todo check user-pass
			//todo server.sendRegistration
			//show success message or error
			return m, gotoRootMenuCmd
		}

	case form.CancelMsg:
		return m, gotoRootMenuCmd
	}

	if m.current == nil {
		return m, nil
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
