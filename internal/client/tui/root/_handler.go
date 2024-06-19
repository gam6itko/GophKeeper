package root

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/gam6itko/goph-keeper/internal/client/tui/loading"
)

type LoginHandler struct {
	prev   tea.Model
	next   tea.Model
	server server.ILoginServer
}

func NewLoginHandler(prev tea.Model, next tea.Model, server server.ILoginServer) *LoginHandler {
	return &LoginHandler{
		prev:   prev,
		next:   next,
		server: server,
	}
}

func (ths LoginHandler) Handle(username string, password string) (tea.Model, tea.Cmd) {
	ctx, fnCancel := context.WithCancel(context.Background())

	s := loading.NewScreen(
		func() tea.Msg {
			if err := ths.server.Login(ctx, username, password); err != nil {
				return loading.NewDoneMsg(
					fmt.Sprintf("ERROR. %s", err),
					func() (tea.Model, tea.Cmd) {
						return ths.prev, nil
					},
				)
			}

			return loading.NewDoneMsg(
				"Login success.",
				func() (tea.Model, tea.Cmd) {
					return ths.next, ths.next.Init()
				},
			)
		},
		fnCancel,
	)

	return s, s.Init()
}

type RegistrationHandler struct {
	prev   tea.Model
	next   tea.Model
	server server.IRegistrationServer
}

func NewRegistrationHandler(prev tea.Model, next tea.Model, server server.IRegistrationServer) *RegistrationHandler {
	return &RegistrationHandler{
		prev:   prev,
		next:   next,
		server: server,
	}
}

func (ths RegistrationHandler) Handle(username string, password string) (tea.Model, tea.Cmd) {
	ctx, fnCancel := context.WithCancel(context.Background())

	s := loading.NewScreen(
		func() tea.Msg {
			if err := ths.server.Register(ctx, username, password); err != nil {
				return loading.NewDoneMsg(
					fmt.Sprintf("ERROR. %s", err),
					func() (tea.Model, tea.Cmd) {
						return ths.prev, nil
					},
				)
			}

			return loading.NewDoneMsg(
				"Registration success. You can login now.",
				func() (tea.Model, tea.Cmd) {
					return ths.next, ths.next.Init()
				},
			)
		},
		fnCancel,
	)

	return s, s.Init()
}
