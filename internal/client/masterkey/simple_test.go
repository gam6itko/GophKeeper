package masterkey_test

import (
	"github.com/gam6itko/goph-keeper/internal/client/masterkey"
	"github.com/stretchr/testify/require"
	"testing"
)

var _ masterkey.IStorage = (*masterkey.MemGuardStorage)(nil)

func TestSimpleStorage(t *testing.T) {
	t.Run("store and load", func(t *testing.T) {
		s := masterkey.NewSimpleStorage()
		require.False(t, s.Has())

		b, err := s.Load()
		require.Nil(t, b)
		require.Error(t, err)

		err = s.Store([]byte("secret key"))
		require.NoError(t, err)
		require.True(t, s.Has())

		b, err = s.Load()
		require.NoError(t, err)
		require.Equal(t, "secret key", string(b))

		s.Clear()
		require.False(t, s.Has())
	})
}
