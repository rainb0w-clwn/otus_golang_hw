package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var s = "abc"

func TestReadDir(t *testing.T) {
	removeAndCloseTempFile := func(f *os.File) {
		t.Helper()
		err := os.Remove(f.Name())
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
	}
	createTempFile := func(t *testing.T, dir string, pattern string) (*os.File, string) {
		t.Helper()
		f, err := os.CreateTemp(dir, pattern)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			removeAndCloseTempFile(f)
		})
		return f, f.Name()
	}
	removeTempDir := func(t *testing.T, dir string) {
		t.Helper()
		err := os.Remove(dir)
		require.NoError(t, err)
	}
	createTempDir := func(t *testing.T, dir string) string {
		t.Helper()
		d, err := os.MkdirTemp(dir, "test")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			removeTempDir(t, d)
		})
		return d
	}
	t.Run("not dir", func(t *testing.T) {
		_, dir := createTempFile(t, "", "")
		_, err := ReadDir(dir)
		require.Error(t, err)
	})
	t.Run("empty dir | subdir", func(t *testing.T) {
		dir := createTempDir(t, "")
		_ = createTempDir(t, dir)
		res, err := ReadDir(dir)
		require.NoError(t, err)
		require.Empty(t, res)
	})
	t.Run("= in name", func(t *testing.T) {
		dir := createTempDir(t, "")
		_, _ = createTempFile(t, dir, "=")
		res, err := ReadDir(dir)
		require.NoError(t, err)
		require.Empty(t, res)
	})
	t.Run("trim space and tabs", func(t *testing.T) {
		dir := createTempDir(t, "")
		f, fN := createTempFile(t, dir, "")
		fN = filepath.Base(fN)
		_, err := io.WriteString(f, s+" \t \t ")
		require.NoError(t, err)
		res, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, s, res[fN].Value)
	})
	t.Run("0x00 to \\n", func(t *testing.T) {
		dir := createTempDir(t, "")
		f, fN := createTempFile(t, dir, "")
		fN = filepath.Base(fN)
		b := string([]byte{0x00})
		_, err := io.WriteString(f, b+s+b)
		require.NoError(t, err)
		res, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, "\n"+s+"\n", res[fN].Value)
	})
	t.Run("remove value | empty value", func(t *testing.T) {
		dir := createTempDir(t, "")
		_, removeFN := createTempFile(t, dir, "")
		removeFN = filepath.Base(removeFN)
		emptyF, emptyFN := createTempFile(t, dir, "")
		emptyFN = filepath.Base(emptyFN)
		_, err := io.WriteString(emptyF, "\n"+s)
		require.NoError(t, err)
		res, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, Environment{
			removeFN: EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			emptyFN: EnvValue{
				Value:      "",
				NeedRemove: false,
			},
		}, res)
	})
}
