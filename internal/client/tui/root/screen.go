package root

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common"
	"github.com/gam6itko/goph-keeper/internal/client/tui/login"
	"github.com/gam6itko/goph-keeper/internal/client/tui/private/root"
)

var buildVersion = "0.0.0"
var buildDate = "never"

const (
	loginOption = iota
	registrationOption
	exitOption
)

type Screen struct {
	list        list.Model
	anotherList tea.Model
	another     bool
}

func NewScreen() *Screen {
	items := []list.Item{
		common.NewListItem("Login", "into account", loginOption),
		common.NewListItem("Registration", "for a new user.", registrationOption),
		common.NewListItem("Exit", "", exitOption),
	}

	m := &Screen{
		list:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		anotherList: root.NewScreen(nil),
	}
	m.list.Title = fmt.Sprintf("GophKeeper. Version: %s. Build: %s.", buildVersion, buildDate)

	return m
}

func (ths Screen) Init() tea.Cmd {
	return nil
}

func (ths Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ths, tea.Quit
		}

		// select option
		if msg.Type == tea.KeyEnter {
			switch ths.list.SelectedItem().(common.ListItem).ID() {
			case loginOption:
				//next := root.NewScreen(ths)
				//return login.NewScreen(ths, NewLoginHandler(ths, next, &server.MockLoginHandler{})), nil
				ths.another = true
				return ths, nil
			case registrationOption:
				return login.NewScreen(ths, NewRegistrationHandler(ths, ths, &server.MockRegistrationHandler{})), nil
			case exitOption:
				return ths, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		ths.list.SetSize(msg.Width, msg.Height)
		ths.anotherList.Update(msg)
	}

	if ths.another {
		return ths.anotherList.Update(msg)
	}

	var cmd tea.Cmd
	ths.list, cmd = ths.list.Update(msg)
	return ths, cmd
}

func (ths Screen) View() string {
	if ths.another {
		return ths.anotherList.View()
	}
	return ths.list.View()
}
