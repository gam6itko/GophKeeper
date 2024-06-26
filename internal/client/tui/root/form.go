package root

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/tui/common/form"
	"github.com/gam6itko/goph-keeper/internal/client/tui/masterkey"
	"regexp"
	"strconv"
	"strings"
)

const (
	RegFormUsernameIndex = iota
	RegFormPasswordIndex
)

const (
	LoginFormUsernameIndex = iota
	LoginFormPasswordIndex
)

var (
	usernameValidate = func(s string) error {
		if !reLogin.MatchString(s) {
			return errors.New("invalid username")
		}
		return nil
	}
)

var reLogin = regexp.MustCompile("^[a-zA-Z0-9]+$")

func newLoginForm() tea.Model {
	username := textinput.New()
	username.Placeholder = "Username"
	username.CharLimit = 32
	username.Focus()
	username.Validate = func(s string) error {
		if !reLogin.MatchString(s) {
			return errors.New("invalid username")
		}
		return nil
	}

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
	username.Validate = usernameValidate

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
func newMasterKeyForm(successRetryCmd tea.Cmd, prev tea.Model) tea.Model {
	input := textinput.New()
	input.Placeholder = "MasterKey"
	input.CharLimit = 32
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'

	inputs := []textinput.Model{
		input,
	}

	f := form.New(inputs, "Enter MasterKey")
	return masterkey.New(f, successRetryCmd, prev)
}

func initCommonStoreInputs() []textinput.Model {
	name := textinput.New()
	name.Placeholder = "Name"
	name.CharLimit = 32
	name.Focus()

	meta := textinput.New()
	meta.Placeholder = "Meta"
	meta.CharLimit = 32

	return []textinput.Model{
		name,
		meta,
	}
}

// newStoreLoginPassForm - форма сохранения логина и пароля.
func newStoreLoginPassForm() tea.Model {
	inputs := initCommonStoreInputs()

	username := textinput.New()
	username.Placeholder = "Username"
	username.CharLimit = 32
	username.Validate = usernameValidate

	password := textinput.New()
	password.Placeholder = "Password"
	password.CharLimit = 32
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	inputs = append(inputs, username, password)

	return form.New(inputs, "Store Login, Password")
}

// newStoreLoginPassForm - форма сохранения текста.
func newStoreTextForm() tea.Model {
	inputs := initCommonStoreInputs()

	text := textinput.New()
	text.Placeholder = "Text"
	text.Width = 100
	//todo support new line

	inputs = append(inputs, text)

	return form.New(inputs, "Store Text")
}

// newStoreLoginPassForm - форма сохранения текста.
func newStoreBinaryForm() tea.Model {
	inputs := initCommonStoreInputs()

	file := textinput.New()
	file.Placeholder = "/file/path"
	file.CharLimit = 100
	file.Width = 100
	file.ShowSuggestions = true

	inputs = append(inputs, file)

	return form.New(inputs, "Store Binary")
}

// See: https://github.com/charmbracelet/bubbletea/blob/master/examples/credit-card-form/main.go
func newStoreBankCardForm() tea.Model {
	inputs := initCommonStoreInputs()

	number := textinput.New()
	number.Placeholder = "4505 **** **** 1234"
	number.CharLimit = 20
	number.Width = 30
	number.Validate = func(s string) error {
		if len(s) > 16+3 {
			return fmt.Errorf("CCN is too long")
		}

		if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
			return fmt.Errorf("CCN is invalid")
		}

		if len(s)%5 == 0 && s[len(s)-1] != ' ' {
			return fmt.Errorf("CCN must separate groups with spaces")
		}

		// The remaining digits should be integers
		c := strings.ReplaceAll(s, " ", "")
		_, err := strconv.ParseInt(c, 10, 64)

		return err
	}

	exp := textinput.New()
	exp.Placeholder = "MM/YY "
	exp.CharLimit = 5
	exp.Width = 5
	exp.Validate = func(s string) error {
		e := strings.ReplaceAll(s, "/", "")
		_, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return fmt.Errorf("EXP is invalid")
		}

		if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
			return fmt.Errorf("EXP is invalid")
		}

		return nil
	}

	cvv := textinput.New()
	cvv.Placeholder = "XXX"
	cvv.CharLimit = 3
	cvv.Width = 5
	cvv.Validate = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	inputs = append(inputs, number, exp, cvv)

	return form.New(inputs, "Store BankCard")
}
