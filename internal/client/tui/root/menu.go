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

	// gotoPrivateMenuMsg нужно вернуться в ЛК.
	gotoPrivateMenuMsg struct{}
	// privateListRequestMsg отобразить список всех сохранённых данных.
	privateListRequestMsg  struct{}
	privateListResponseMsg struct {
		list []server.PrivateDataListItemDTO
	}
	privateDataLoadMsg struct {
		id uint32
	}

	// privateLogoutMsg выйти из ЛК.
	privateLogoutMsg struct{}
	// privateStoreStartMsg - начать процедуру ввода новых данных.
	privateStoreStartMsg struct {
		dataType server.TPrivateData
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
				return privateListRequestMsg{}
			}),
			common.NewCmdItem("Store", "login password", func() tea.Msg {
				return privateStoreStartMsg{dataType: server.TypeLoginPass}
			}),
			common.NewCmdItem("Store", "text", func() tea.Msg {
				return privateStoreStartMsg{dataType: server.TypeText}
			}),
			common.NewCmdItem("Store", "binary", func() tea.Msg {
				return privateStoreStartMsg{dataType: server.TypeBinary}
			}),
			common.NewCmdItem("Store", "bank card", func() tea.Msg {
				return privateStoreStartMsg{dataType: server.TypeBankCard}
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

// newPrivateDataList создаёт меню со всеми имеющимися данными на сервере.
// Пользователь может выбрать какие данные нужно загрузить и отобразить.
func newPrivateDataList(title string, width, height int, privateList []server.PrivateDataListItemDTO) tea.Model {
	items := make([]list.Item, len(privateList), len(privateList)+1)
	for i, item := range privateList {
		id := item.ID
		items[i] = common.NewCmdItem(item.Name, item.Type.String(), func() tea.Msg {
			return privateDataLoadMsg{id: id}
		})
	}
	items = append(
		items,
		common.NewCmdItem("<- Go Back", "Go to prev menu", func() tea.Msg {
			return gotoPrivateMenuMsg{}
		}),
	)

	l := list.New(
		items,
		list.NewDefaultDelegate(),
		width,
		height,
	)
	l.Title = title

	return cmd.New(l)
}
