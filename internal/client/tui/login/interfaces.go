package login

import tea "github.com/charmbracelet/bubbletea"

type ISubmitHandler interface {
	Handle(username string, password string) (tea.Model, tea.Cmd)
}
