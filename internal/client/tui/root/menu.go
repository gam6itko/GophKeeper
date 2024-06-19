package root

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/menu/cmd"
	"github.com/gam6itko/goph-keeper/internal/client/tui/menu/common"
)

type (
	gotoRootMenuMsg struct{}
	// gotoLoginMsg управляющая модель должна отобразить форму аутентификации.
	gotoLoginMsg struct{}
	// gotoRegistrationMsg управляющая модель должна отобразить форму регистрации.
	gotoRegistrationMsg struct{}

	// privateListMsg отобразить список всех сохранённых данных.
	privateListMsg struct{}
	// privateLogoutMsg выйти из ЛК.
	privateLogoutMsg struct{}
	// Начать процедуру ввода новых данных.
	privateStoreMsg struct {
		t server.PrivateDataType
	}
)

func gotoRootMenuCmd() tea.Msg {
	return gotoRootMenuMsg{}
}
func gotoLoginCmd() tea.Msg {
	return gotoLoginMsg{}
}
func gotoRegistrationCmd() tea.Msg {
	return gotoRegistrationMsg{}
}

// newRootMenu создаёт главное меню.
func newRootMenu(title string, width, height int) tea.Model {
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

	return cmd.New(l)
}

// newPrivateMenu создаёт меню ЛК пользователя.
func newPrivateMenu(title string, width, height int) tea.Model {
	l := list.New(
		[]list.Item{
			common.NewCmdItem("List", "of stored entries", func() tea.Msg {
				return privateListMsg{}
			}),
			common.NewCmdItem("Store", "login password", func() tea.Msg {
				return privateStoreMsg{t: server.TypeLoginPass}
			}),
			common.NewCmdItem("Store", "text", func() tea.Msg {
				return privateStoreMsg{t: server.TypeText}
			}),
			common.NewCmdItem("Store", "binary", func() tea.Msg {
				return privateStoreMsg{t: server.TypeBinary}
			}),
			common.NewCmdItem("Store", "bank card", func() tea.Msg {
				return privateStoreMsg{t: server.TypeBankCard}
			}),
			common.NewCmdItem("Logout", "", func() tea.Msg {
				return privateLogoutMsg{}
			}),
		},
		list.NewDefaultDelegate(),
		width,
		height,
	)
	l.Title = title

	return cmd.New(l)
}
