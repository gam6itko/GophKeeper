package serialize

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestText_Serialize(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	s := Text{}
	b, err := s.Serialize(text)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	text2, err := s.Deserialize(b)
	require.NoError(t, err)
	require.Equal(t, text, text2)
}
