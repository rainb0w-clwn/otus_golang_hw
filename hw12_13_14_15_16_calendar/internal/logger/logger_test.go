package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const msg = "log msg"
	t.Run("log with exact level", func(t *testing.T) {
		out := &bytes.Buffer{}
		logg := New(Warning, out)
		logg.Warning(msg)
		require.Contains(t, out.String(), msg)
	})
	t.Run("log with higher level", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Info, out)
		log.Warning(msg)
		require.Contains(t, out.String(), msg)
	})
	t.Run("log with lower level", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Error, out)
		log.Warning(msg)
		require.Empty(t, out.String())
	})
}
