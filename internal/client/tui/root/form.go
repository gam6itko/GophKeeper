package root

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/masterkey"
)

const (
	RegFormUsernameIndex = iota
	RegFormPasswordIndex
)

const (
	LoginFormUsernameIndex = iota
	LoginFormPasswordIndex
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

// newMasterKeyForm создать форму для ввода мастер-ключа, которым зашифрованы все данные.
func newMasterKeyForm(successRetryMsg tea.Msg, prev tea.Model) tea.Model {
	input := textinput.New()
	input.Placeholder = "MasterKey"
	input.CharLimit = 32
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'

	inputs := []textinput.Model{
		input,
	}

	f := form.New(inputs, "Enter MasterKey")
	return masterkey.New(f, successRetryMsg, prev)
}
