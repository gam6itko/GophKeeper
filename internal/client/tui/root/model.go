package root

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/loading"
	masterkey_form "github.com/gam6itko/goph-keeper/internal/client/tui/masterkey"
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
	stateOnPrivateMenu
	stateOnRegistrationForm
	stateStorePrivateData
	stateLoadPrivate
)

type Model struct {
	server  server.IServer
	storage masterkey.IStorage

	width, height int
	current       tea.Model
	prev          tea.Model
	state         appState
	cancelFunc    *context.CancelFunc
}

func New(server server.IServer, storage masterkey.IStorage) *Model {
	return &Model{
		state:   stateIdle,
		current: newRootMenu(fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate), 0, 0),

		server:  server,
		storage: storage,
	}
}

func (m Model) Init() tea.Cmd {
	return gotoRootMenuCmd
}

// Update - единая точка обработки всех команд от моделей.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" {
			if m.cancelFunc != nil {
				fn := *m.cancelFunc
				fn()
				m.cancelFunc = nil
				return m, nil
			}
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

	// Пользователь заполнил форму и нажал Submit...
	case form.SubmitMsg:
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		switch m.state {
		// ... на форме входа в систему.
		case stateOnLoginFrom:
			prev := m.current
			m.current = loading.New(
				func() tea.Msg {
					err := m.server.Login(
						ctx,
						server.LoginDTO{
							Username: msg.Values[LoginFormUsernameIndex],
							Password: msg.Values[LoginFormPasswordIndex],
						},
					)
					if err == nil {
						return loading.ResponseResultMsg{
							Success: true,
							Message: "Login successful",
						}
					} else {
						return loading.ResponseResultMsg{
							Success: false,
							Message: fmt.Sprintf("Error. %s", err),
						}
					}
				},
				func() tea.Msg {
					return loading.SuccessMsg{}
				},
				func() tea.Msg {
					return loading.FailMsg{GoToModel: prev}
				},
			)
			return m, m.current.Init()
		// ... на форме регистрации.
		case stateOnRegistrationForm:
			log.Printf("submit reg")
			//todo check user-pass
			//todo server.sendRegistration
			//show success message or error
			return m, gotoRootMenuCmd
		}

	// Запрос к серверу прошёл успешно.
	case loading.SuccessMsg:
		switch m.state {
		case stateOnLoginFrom:
			m.state = stateOnPrivateMenu
			m.cancelFunc = nil
			m.current = newPrivateMenu("Private menu", m.width, m.height)
			return m, m.current.Init()

		case stateOnRegistrationForm:
			//return m, m.current.Init()
		}

	// Запрос к серверу завершился ошибкой.
	case loading.FailMsg:
		model := msg.GoToModel
		m.current = model
		m.cancelFunc = nil
		return m, model.Init()

	// Отмена ввода в форме.
	case form.CancelMsg:
		return m, gotoRootMenuCmd

	//// private menu msg

	case gotoPrivateMenuMsg:
		m.current = newPrivateMenu("Private menu", m.width, m.height)
		return m, m.current.Init()
	// Пользователь захотел посмотреть список данных имеющихся на сервере.
	case privateListRequestMsg:
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		m.current = loading.New(
			func() tea.Msg {
				list, err := m.server.List(ctx)
				if err == nil {
					return privateListResponseMsg{
						list: list,
					}
				} else {
					return loading.ResponseResultMsg{
						Success: false,
						Message: fmt.Sprintf("Error. %s", err),
					}
				}
			},
			gotoRootMenuCmd,
			func() tea.Msg {
				return loading.FailMsg{GoToModel: m.current}
			},
		)
		return m, m.current.Init()

	// Пришёл ответ со списком всех имеющихся данных на сервере.
	case privateListResponseMsg:
		m.current = newPrivateDataList("Your private date. Click to view", m.width, m.height, msg.list)
		return m, m.current.Init()

	case privateDataLoadMsg:
		if !m.storage.Has() {
			m.current = newMasterKeyForm(msg, m.current)
			return m, m.current.Init()
		}

		// Начинаем грузить с сервера приватные данные
		//ctx, cancelFunc := context.WithCancel(context.Background())
		//m.cancelFunc = &cancelFunc

		//todo pd, err := m.server.Load(ctx, msg.id)
		//todo loading
		//	on success: decode data, show type form

		return m, m.current.Init()

	// Пользователь нажал Logout.
	case privateLogoutMsg:
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		prev := m.current
		m.current = loading.New(
			func() tea.Msg {
				err := m.server.Logout(ctx)
				if err == nil {
					return loading.ResponseResultMsg{
						Success: true,
						Message: "Logout successful",
					}
				} else {
					return loading.ResponseResultMsg{
						Success: false,
						Message: fmt.Sprintf("Error. %s", err),
					}
				}
			},
			gotoRootMenuCmd,
			func() tea.Msg {
				return loading.FailMsg{GoToModel: prev}
			},
		)
		return m, m.current.Init()

	case masterkey_form.SubmitMsg:
		if err := m.storage.Store(msg.Key); err != nil {
			log.Printf("store error: %s", err)
			//todo show error message
			return m, nil
		}
		return m, func() tea.Msg {
			return msg.RetryMsg
		}
	case masterkey_form.CancelMsg:
	}

	//if m.current == nil {
	//	return m, nil
	//}

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
