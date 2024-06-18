package root

//
//import (
//	"fmt"
//	"github.com/charmbracelet/bubbles/list"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
//)
//
//const (
//	listOption = iota
//	getOption
//	storeLPOption
//	storeTextOption
//	storeBinaryOption
//	storeCardOption
//	logoutOption
//)
//
//// Screen shows private user area after login.
//type Screen struct {
//	prev tea.Model
//	list list.Model
//}
//
//func NewScreen(prev tea.Model) *Screen {
//	items := []list.Item{
//		common.NewListItem("List", "of stored items", listOption),
//		//common.NewListItem("Get", "item by ID", getOption),
//		//common.NewListItem("Store", "login password", storeLPOption),
//		//common.NewListItem("Store", "text", storeTextOption),
//		//common.NewListItem("Store", "binary", storeBinaryOption),
//		//common.NewListItem("Store", "bank card", storeCardOption),
//		common.NewListItem("Logout", "", logoutOption),
//	}
//
//	m := &Screen{
//		prev: prev,
//		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
//	}
//	m.list.Title = fmt.Sprintf("User private area.")
//
//	return m
//}
//
//func (ths Screen) Init() tea.Cmd {
//	return nil
//}
//
//func (ths Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		if msg.String() == "ctrl+c" {
//			return ths, tea.Quit
//		}
//		if msg.String() == "esc" {
//			return ths.prev, nil
//		}
//
//		if msg.Type == tea.KeyEnter {
//			switch ths.list.SelectedItem().(common.ListItem).ID() {
//			case listOption:
//				//todo server.get list
//			case getOption:
//				// id enter
//				//master pass enter
//			case storeLPOption:
//				// screen login pass
//				// screen master pass
//				// send Store(type, []byte)
//			case logoutOption:
//				//logout
//				return ths.prev, nil
//			}
//		}
//
//	case tea.WindowSizeMsg:
//		ths.list.SetSize(msg.Width, msg.Height)
//	}
//
//	var cmd tea.Cmd
//	ths.list, cmd = ths.list.Update(msg)
//	return ths, cmd
//}
//
//func (ths Screen) View() string {
//	return ths.list.View()
//}
