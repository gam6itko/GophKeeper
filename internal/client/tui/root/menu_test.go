package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gam6itko/goph-keeper/internal/client/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMenuBuildSmoke(t *testing.T) {
	t.Run("newRootMenu", func(t *testing.T) {
		m := newRootMenu("test", 10, 10)
		assert.Implements(t, (*tea.Model)(nil), m)
	})

	t.Run("newPrivateMenu", func(t *testing.T) {
		m := newPrivateMenu("test", 10, 10)
		assert.Implements(t, (*tea.Model)(nil), m)
	})

	t.Run("newPrivateDataList", func(t *testing.T) {
		m := newPrivateDataList("test", 10, 10, []server.PrivateDataListItemDTO{})
		assert.Implements(t, (*tea.Model)(nil), m)
	})
}
