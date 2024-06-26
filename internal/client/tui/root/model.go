package root

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"github.com/gam6itko/goph-keeper/internal/client/private"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/errmsg"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/info"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/loading"
	masterkey_form "github.com/gam6itko/goph-keeper/internal/client/tui/masterkey"
	"log"
	"os"
)

var buildVersion = "0.0.0"
var buildDate = "never"

// masterKeyCheckSign - строка которая должна быть корректно в расшифрованна местер-ключом сообщении.
// Если строки не идентичны, то мастер-ключ был указан неверно.
var masterKeyCheckSign = "WTF"

type (
	// gotoModelMsg нужно отобразить на экране другую модель.
	gotoModelMsg struct {
		// model которая будет отображена на экране.
		model tea.Model
		// newState если не stateUndefined то Model.state будет изменён на то что указано.
		newState appState
	}

	// showErrorMsg - нужно отобразить сообщение об ошибке.
	showErrorMsg struct {
		// err ошибка
		err error
		// cmd которую нужно выполнить когда пользователь прочитал ошибку и решил идти дальше.
		cmd tea.Cmd
	}

	showPrivateDataMsg struct {
		dto *server.PrivateDataDTO
	}
)

var (
	// newCmd обертка чтобы отправить msg.
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

	newShowErrCmd = func(err error, cmd tea.Cmd) tea.Cmd {
		return func() tea.Msg {
			return showErrorMsg{
				err: err,
				cmd: cmd,
			}
		}
	}
)

type appState int

const (
	// stateUndefined - состояние не указано.
	// Нужно чтобы newGotoModelCmd не менять текущее состояние.
	stateUndefined appState = iota
	stateStartup
	stateOnRootMenu
	stateOnLoginFrom
	stateOnPrivateMenu
	stateOnRegistrationForm
	stateStoreLoginPass
	stateStoreText
	stateStoreBinary
	stateStoreBankCard
)

// Model - самая главная модель этой программы.
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
		m.cancelFunc = nil
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

	// Велено отобразить ошибку и отправить команду при закрытии.
	case showErrorMsg:
		m.cancelFunc = nil
		m.current = errmsg.New(msg.err, msg.cmd)
		return m, nil

	// Велено идти в ЛК пользователя.
	case gotoPrivateMenuMsg:
		m.cancelFunc = nil
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
						err: err,
						cmd: newGotoModelCmd(gotoModelFail, 0),
					}
				}
			},
		)
		return m, m.current.Init()

	// Пришёл ответ со списком всех имеющихся данных на сервере.
	case privateListResponseMsg:
		m.cancelFunc = nil
		m.current = newPrivateDataList("Your private date. Click to view", m.width, m.height, msg.list)
		return m, m.current.Init()

	// Пользователь запросил данные с сервера. Отправляем запрос-ждм ответа.
	case privateDataLoadMsg:
		// Проверяем есть ли мастер-ключ. Если нет, то просим ввести.
		if !m.storage.Has() {
			m.current = newMasterKeyForm(newCmd(msg), m.current)
			return m, m.current.Init()
		}

		// Начинаем грузить с сервера приватные данные.
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancelFunc = &cancelFunc

		dataID := msg.id
		m.current = loading.New(
			func() tea.Msg {
				dto, err := m.server.Load(ctx, dataID)
				if err == nil {
					return showPrivateDataMsg{dto}
				}
				return gotoPrivateMenuMsg{}
			},
		)
		return m, m.current.Init()

	// Данные пришли с сервера. Декодируем и отображаем.
	case showPrivateDataMsg:
		m.cancelFunc = nil

		data := m.decodeData(msg.dto.Data)
		buff := bytes.NewBuffer(data[3:])
		decoder := gob.NewDecoder(buff)

		switch msg.dto.Type {
		case server.TypeLoginPass:
			dto := server.LoginPassDTO{}
			err := decoder.Decode(&dto)
			if err != nil {
				log.Fatalf("login pass decode error: %s", err)
			}
			if dto.Sign != masterKeyCheckSign {
				m.storage.Clear()
				m.current = errmsg.New(errors.New("invalid master key"), gotoPrivateMenuCmd)
				return m, nil
			}
			m.current = info.NewModel(
				map[string]string{
					"login   ": dto.Login,
					"password": dto.Password,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil

		case server.TypeText:
			dto := server.TextDTO{}
			err := decoder.Decode(&dto)
			if err != nil {
				log.Fatalf("text decode error: %s", err)
			}
			if dto.Sign != masterKeyCheckSign {
				m.storage.Clear()
				m.current = errmsg.New(errors.New("invalid master key"), gotoPrivateMenuCmd)
				return m, nil
			}
			m.current = info.NewModel(
				map[string]string{
					"text": dto.Text,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil

		// Скачанные двоичные данные просто сохраняем в файл.
		//todo дать пользователю возможность ввести название файла куда сор+хранить.
		case server.TypeBinary:
			dto := server.BinaryDTO{}
			err := decoder.Decode(&dto)
			if err != nil {
				log.Fatalf("binary decode fail. %s", err)
			}
			if dto.Sign != masterKeyCheckSign {
				m.storage.Clear()
				m.current = errmsg.New(errors.New("invalid master key"), gotoPrivateMenuCmd)
				return m, nil
			}

			f, err := os.CreateTemp("/tmp", "*.bin")
			_, err = f.Write(dto.Data)
			if err != nil {
				log.Fatal(err)
			}

			m.current = info.NewModel(
				map[string]string{
					"save to file": f.Name(),
				},
				gotoPrivateMenuCmd,
			)
			return m, nil

		case server.TypeBankCard:
			dto := server.BankCardDTO{}
			err := decoder.Decode(&dto)
			if err != nil {
				log.Fatalf("bank card decode error: %s", err)
			}
			m.current = info.NewModel(
				map[string]string{
					"number ": dto.Number,
					"expires": dto.Expires,
					"cvv    ": dto.CVV,
				},
				gotoPrivateMenuCmd,
			)
			return m, nil
		}

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
					return gotoRootMenuMsg{}
				}
				return showErrorMsg{
					err: err,
					cmd: newGotoModelCmd(gotoModelFail, 0),
				}
			},
		)
		return m, m.current.Init()

	case masterkey_form.SubmitMsg:
		key := m.crypt.KeyFit(msg.Key)
		if err := m.storage.Store(key); err != nil {
			// Если при сохранении мастер-ключа произошла ошибка, то возвращаемся в ЛК.
			m.current = errmsg.New(err, gotoPrivateMenuCmd)
			return m, nil
		}

		return m, msg.RetryCmd

	case masterkey_form.CancelMsg:
		return m, gotoPrivateMenuCmd

	case privateStoreStartMsg:
		switch msg.dataType {
		case server.TypeLoginPass:
			return m, newGotoModelCmd(newStoreLoginPassForm(), stateStoreLoginPass)
		case server.TypeText:
			return m, newGotoModelCmd(newStoreTextForm(), stateStoreText)
		case server.TypeBinary:
			return m, newGotoModelCmd(newStoreBinaryForm(), stateStoreBinary)
		case server.TypeBankCard:
			return m, newGotoModelCmd(newStoreBankCardForm(), stateStoreBankCard)
		}
		return m, nil
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
						return gotoPrivateMenuMsg{}
					}
					return showErrorMsg{
						err: err,
						cmd: newGotoModelCmd(gotoModelFail, 0),
					}
				},
			)
			return m, m.current.Init()

		// Отмена ввода в форме.
		case form.CancelMsg:
			return m, gotoRootMenuCmd
		}

	case stateOnRegistrationForm:
		switch msg := msg.(type) {
		case form.SubmitMsg:
			ctx, cancelFunc := context.WithCancel(context.Background())
			m.cancelFunc = &cancelFunc
			gotoModelFail := m.current
			m.current = loading.New(
				func() tea.Msg {
					err := m.server.Registration(
						ctx,
						server.RegistrationDTO{
							Username: msg.Values[RegFormUsernameIndex],
							Password: msg.Values[RegFormPasswordIndex],
						},
					)
					if err == nil {
						//todo отображать сообщение можно покрасивее.
						model := info.NewModel(
							map[string]string{
								"message": "Registration success. You can login now",
							},
							gotoRootMenuCmd,
						)
						return gotoModelMsg{
							model: model,
						}
					}
					return showErrorMsg{
						err: err,
						cmd: newGotoModelCmd(gotoModelFail, 0),
					}
				},
			)
			return m, m.current.Init()
		// Отмена ввода в форме.
		case form.CancelMsg:
			return m, gotoRootMenuCmd
		}

	case stateStoreLoginPass:
		switch msg := msg.(type) {
		case form.SubmitMsg:
			idxName, idxMeta, idxLogin, idxPass := 0, 1, 2, 3
			dto := server.LoginPassDTO{
				Sign:     masterKeyCheckSign,
				Login:    msg.Values[idxLogin],
				Password: msg.Values[idxPass],
			}
			return m.storeData(msg.Values[idxName], msg.Values[idxMeta], dto)

		case form.CancelMsg:
			return m, gotoPrivateMenuCmd
		}

	case stateStoreText:
		switch msg := msg.(type) {
		case form.SubmitMsg:
			idxName, idxMeta, idxText := 0, 1, 2
			dto := server.TextDTO{
				Sign: masterKeyCheckSign,
				Text: msg.Values[idxText],
			}
			return m.storeData(msg.Values[idxName], msg.Values[idxMeta], dto)

		case form.CancelMsg:
			return m, gotoPrivateMenuCmd
		}

	// Велено на сервер отправить двоичные данные.
	// Бинарные данные берём из файла.
	case stateStoreBinary:
		switch msg := msg.(type) {
		case form.SubmitMsg:
			idxName, idxMeta, idxFilePath := 0, 1, 2
			data, err := os.ReadFile(msg.Values[idxFilePath])
			if err != nil {
				return m, newShowErrCmd(err, gotoPrivateMenuCmd)
			}

			dto := server.BinaryDTO{
				Sign: masterKeyCheckSign,
				Data: data,
			}
			return m.storeData(msg.Values[idxName], msg.Values[idxMeta], dto)

		case form.CancelMsg:
			return m, gotoPrivateMenuCmd
		}

	case stateStoreBankCard:
		switch msg := msg.(type) {
		case form.SubmitMsg:
			idxName, idxMeta, idxNumber, idxExpires, idxCVV := 0, 1, 2, 3, 4
			dto := server.BankCardDTO{
				Sign:    masterKeyCheckSign,
				Number:  msg.Values[idxNumber],
				Expires: msg.Values[idxExpires],
				CVV:     msg.Values[idxCVV],
			}

			if err := dto.Validate(); err != nil {
				m.current = errmsg.New(err, gotoPrivateMenuCmd)
				return m, m.current.Init()
			}

			return m.storeData(msg.Values[idxName], msg.Values[idxMeta], dto)

		case form.CancelMsg:
			return m, gotoPrivateMenuCmd
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

func (m Model) storeData(name string, meta string, dto any) (tea.Model, tea.Cmd) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	m.cancelFunc = &cancelFunc

	gotoModelFail := m.current
	m.current = loading.New(
		func() tea.Msg {
			b := []byte(nil)
			buff := bytes.NewBuffer(b)
			encoder := gob.NewEncoder(buff)
			err := encoder.Encode(dto)
			if err != nil {
				log.Fatalf("login pass encode error: %s", err)
			}

			data := server.PrivateDataDTO{
				BasePrivateDataDTO: server.BasePrivateDataDTO{
					Type: private.TypeLoginPass,
					Name: name,
					Meta: meta,
				},
				Data: b,
			}
			err = m.server.Store(ctx, data)
			if err == nil {
				//todo отображать сообщение можно покрасивее.
				return gotoModelMsg{
					model: info.NewModel(
						map[string]string{
							"message": "Store success",
						},
						gotoPrivateMenuCmd,
					),
					newState: stateOnPrivateMenu,
				}
			}

			return showErrorMsg{
				err: err,
				cmd: newGotoModelCmd(gotoModelFail, 0),
			}
		},
	)
	return m, m.current.Init()
}
