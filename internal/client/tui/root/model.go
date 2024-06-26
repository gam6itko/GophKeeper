package root

import (
	"bytes"
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"github.com/gam6itko/goph-keeper/internal/client/serialize"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/errmsg"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/info"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/loading"
	masterkey_form "github.com/gam6itko/goph-keeper/internal/client/tui/masterkey"
	"log"
)

var buildVersion = "0.0.0"
var buildDate = "never"

// masterKeyCheckSign - 3 байта которые должны быть в расшифрованном местер-ключом сообщении.
// Если данные не начинаются на эти 3 байта, то мастер ключ не верен.
var masterKeyCheckSign = []byte("WTF")

type (
	gotoModelMsg struct {
		// model которая будет отображена на экране.
		model tea.Model
		// newState если не stateUndefined то Model.state будет изменён на то что указано.
		newState appState
	}

	showErrorMsg struct {
		gotoModel tea.Model
		err       error
	}

	showPrivateDataMsg struct {
		dto *server.PrivateDataDTO
	}
)

var (
	newCmd = func(msg tea.Msg) tea.Cmd {
		return func() tea.Msg {
			return msg
		}
	}
	newGotoModelCmd = func(model tea.Model, state appState) tea.Cmd {
		return func() tea.Msg {
			return gotoModelMsg{
				model:    model,
				newState: state,
			}
		}
	}
	gotoPrivateMenuCmd = func() tea.Msg {
		return gotoPrivateMenuMsg{}
	}
)

type appState int

const (
	stateUndefined appState = iota
	stateStartup
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
	crypt   encrypt.ICrypt

	width, height int
	current       tea.Model
	prev          tea.Model
	state         appState
	cancelFunc    *context.CancelFunc
}

func New(server server.IServer, storage masterkey.IStorage, crypt encrypt.ICrypt) *Model {
	return &Model{
		state:   stateStartup,
		current: newRootMenu(fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate), 0, 0),

		server:  server,
		storage: storage,
		crypt:   crypt,
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

	case gotoModelMsg:
		if msg.newState != stateUndefined {
			m.state = msg.newState
		}
		m.current = msg.model
		return m, m.current.Init()

	case gotoRootMenuMsg:
		return m, newGotoModelCmd(
			newRootMenu(
				fmt.Sprintf("GophKeeper. Version: %s. Build: %s", buildVersion, buildDate),
				m.width,
				m.height,
			),
			stateOnRootMenu,
		)

	case gotoLoginMsg:
		return m, newGotoModelCmd(newLoginForm(), stateOnLoginFrom)

	case gotoRegistrationMsg:
		return m, newGotoModelCmd(newRegistrationForm(), stateOnRegistrationForm)

	case showErrorMsg:
		m.current = errmsg.New(msg.err, newGotoModelCmd(msg.gotoModel, 0))
		return m, nil

	case gotoPrivateMenuMsg:
		return m, newGotoModelCmd(
			newPrivateMenu("Private menu", m.width, m.height),
			stateOnPrivateMenu,
		)

	// Пользователь захотел посмотреть список данных имеющихся на сервере.
	case privateListRequestMsg:
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		gotoModelFail := m.current
		m.current = loading.New(
			func() tea.Msg {
				list, err := m.server.List(ctx)
				if err == nil {
					return privateListResponseMsg{
						list: list,
					}
				} else {
					return showErrorMsg{
						err:       err,
						gotoModel: gotoModelFail,
					}
				}
			},
		)
		return m, m.current.Init()

	// Пришёл ответ со списком всех имеющихся данных на сервере.
	case privateListResponseMsg:
		m.current = newPrivateDataList("Your private date. Click to view", m.width, m.height, msg.list)
		return m, m.current.Init()

	case privateDataLoadMsg:
		// Проверяем есть ли мастер-ключ. Если нет, то просим ввести.
		if !m.storage.Has() {
			m.current = newMasterKeyForm(newCmd(msg), m.current)
			return m, m.current.Init()
		}

		// Начинаем грузить с сервера приватные данные
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		dataID := msg.id
		m.current = loading.New(
			func() tea.Msg {
				dto, err := m.server.Load(ctx, dataID)
				if err == nil {
					return loading.DoneCmd{
						Cmd: func() tea.Msg {
							return showPrivateDataMsg{dto}
						},
					}
				} else {
					return loading.DoneCmd{
						Cmd: gotoPrivateMenuCmd,
					}
				}
			},
		)
		return m, m.current.Init()

	case showPrivateDataMsg:
		data := m.decodeData(msg.dto.Data)
		sign := data[:3]
		// Если мастер-ключ верен, то префикс должен совпасть.
		if !bytes.Equal(sign, masterKeyCheckSign) {
			log.Fatal("invalid master key signature")
		}

		data = data[3:]

		switch msg.dto.Type {
		case server.TypeLoginPass:
			ser := serialize.LoginPass{}
			dto, err := ser.Deserialize(data)
			if err != nil {
				log.Fatal("invalid login pass")
			}
			m.current = info.NewModel(
				map[string]string{
					"login":    dto.Login,
					"password": dto.Password,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil

		case server.TypeText:
			ser := serialize.Text{}
			text, err := ser.Deserialize(data)
			if err != nil {
				log.Fatal("invalid text")
			}
			m.current = info.NewModel(
				map[string]string{
					"text": text,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil

		case server.TypeBinary:
		case server.TypeBankCard:
			ser := serialize.BankCard{}
			dto, err := ser.Deserialize(data)
			if err != nil {
				log.Fatal("invalid text")
			}
			m.current = info.NewModel(
				map[string]string{
					"number":  dto.Number,
					"expires": dto.Expires,
					"cvv":     dto.CVV,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil
		}
		//todo type switch
		//  data decode
		//	masterkey-ckecksum validate
		//  gotoModel - type view

		return m, nil

	// Пользователь нажал Logout.
	case privateLogoutMsg:
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		gotoModelFail := m.current
		m.current = loading.New(
			func() tea.Msg {
				err := m.server.Logout(ctx)
				if err == nil {
					return loading.DoneCmd{
						Cmd: gotoRootMenuCmd,
					}
				} else {
					return loading.DoneCmd{
						Cmd: func() tea.Msg {
							return showErrorMsg{
								err:       err,
								gotoModel: gotoModelFail,
							}
						},
					}
				}
			},
		)
		return m, m.current.Init()

	case masterkey_form.SubmitMsg:
		key := m.crypt.KeyFit(msg.Key)
		if err := m.storage.Store(key); err != nil {
			// Если при сохранении мастер-ключа произошла ошибка, тогда возвращаемся в ЛК.
			m.current = errmsg.New(err, gotoPrivateMenuCmd)
			return m, nil
		}

		return m, msg.RetryCmd

	case masterkey_form.CancelMsg:
		return m, gotoPrivateMenuCmd

	}

	// Обработка если приложение в конкретном состоянии.
	switch m.state {
	case stateOnLoginFrom:
		switch msg := msg.(type) {
		// Пользователь заполнил форму и нажал Submit на форме входа в систему.
		case form.SubmitMsg:
			ctx, cancelFunc := context.WithCancel(context.Background())
			m.cancelFunc = &cancelFunc
			gotoModelFail := m.current
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
						return loading.DoneCmd{
							Cmd: newGotoModelCmd(
								newPrivateMenu("Private menu", m.width, m.height),
								stateOnPrivateMenu,
							),
						}
					} else {
						return loading.DoneCmd{
							Cmd: func() tea.Msg {
								return showErrorMsg{
									err:       err,
									gotoModel: gotoModelFail,
								}
							},
						}
					}
				},
			)
			return m, m.current.Init()

		// Отмена ввода в форме.
		case form.CancelMsg:
			return m, gotoRootMenuCmd
		}

	case stateOnRegistrationForm:
		switch msg.(type) {
		case form.SubmitMsg:
			//ctx, cancelFunc := context.WithCancel(context.Background())
			//m.cancelFunc = &cancelFunc
			//log.Printf("submit reg")
			//todo check user-pass
			//todo server.sendRegistration
			//show success message or error
			return m, gotoRootMenuCmd
		// Отмена ввода в форме.
		case form.CancelMsg:
			return m, gotoRootMenuCmd
		}
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

func (m Model) decodeData(data []byte) []byte {
	key, err := m.storage.Load()
	if err != nil {
		log.Fatal(err)
	}

	c := encrypt.NewAESCrypt()
	result, err := c.Decrypt(key, data)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
