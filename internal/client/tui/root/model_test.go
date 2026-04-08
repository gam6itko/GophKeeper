package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/gam6itko/goph-keeper/internal/client/masterkey/encrypt"
	"github.com/gam6itko/goph-keeper/internal/client/server/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Update(t *testing.T) {
	s := &mock.Server{}
	memStorage := &masterkey.SimpleStorage{}
	crypt := encrypt.NewAESCrypt()

	m := New(s, memStorage, crypt)

	t.Run("quit", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		_, cmd := m.Update(msg)
		assert.NotEmpty(t, cmd)
		assert.Equal(t, tea.QuitMsg{}, cmd())
	})

	t.Run("set window size", func(t *testing.T) {
		assert.Empty(t, m.width)
		assert.Empty(t, m.height)

		m2, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 200})
		assert.NotEmpty(t, m2)
		assert.Empty(t, cmd)
		// Не меняется изначальная модель, т.к. receiver передаётся по значению.
		assert.Empty(t, m.width)
		assert.Empty(t, m.height)
		// Изменения вносятся в новую модель.
		assert.IsType(t, Model{}, m2)
		assert.Equal(t, 100, m2.(Model).width)
		assert.Equal(t, 200, m2.(Model).height)
	})
}
