package root

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
)

func newLoginForm() tea.Model {
	username := textinput.New()
	username.Placeholder = "Username"
	username.CharLimit = 32
	username.CharLimit = 64
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Password"
	password.CharLimit = 32
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	inputs := []textinput.Model{
		username,
		password,
	}

	return form.New(inputs, "Login")
}

func newRegistrationForm() tea.Model {
	username := textinput.New()
	username.Placeholder = "Username"
	username.CharLimit = 32
	username.CharLimit = 64
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Password"
	password.CharLimit = 32
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	inputs := []textinput.Model{
		username,
		password,
	}

	return form.New(inputs, "Registration")
}
