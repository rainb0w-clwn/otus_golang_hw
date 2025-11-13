package main

import (
	"errors"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func isUnixLike() bool {
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd", "netbsd", "openbsd", "solaris", "android":
		return true
	}
	return false
}

func TestCopy(t *testing.T) {
	removeAndCloseTempFile := func(f *os.File) {
		t.Helper()
		err := os.Remove(f.Name())
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)
	}
	createTempFile := func(t *testing.T) (*os.File, string) {
		t.Helper()
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			removeAndCloseTempFile(f)
		})
		return f, f.Name()
	}
	t.Run(ErrFromFileNotExists.Error(), func(t *testing.T) {
		_, err := Copy("", "", 0, 0)
		require.Truef(t, errors.Is(err, ErrFromFileNotExists), "actual err - %v", err)
	})
	t.Run(ErrFromFileOpen.Error(), func(t *testing.T) {
		f, err := os.OpenFile(".tmp", os.O_RDWR|os.O_CREATE, os.ModeExclusive)
		require.NoError(t, err)
		defer removeAndCloseTempFile(f)
		_, err = Copy(f.Name(), "", 0, 0)
		require.Truef(t, errors.Is(err, ErrFromFileOpen), "actual err - %v", err)
	})
	t.Run("fromFile is dir", func(t *testing.T) {
		fromFile, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			require.NoError(t, err)
		}(fromFile)
		_, err = Copy(fromFile, "", 0, 0)
		require.Truef(t, errors.Is(err, ErrFromFileUnsupported), "actual err - %v", err)
	})
	t.Run(ErrUnsupportedOffsetLimit.Error(), func(t *testing.T) {
		_, fromPath := createTempFile(t)
		_, toPath := createTempFile(t)
		_, err := Copy(fromPath, toPath, -1, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedOffsetLimit), "actual err - %v", err)
		_, err = Copy(fromPath, toPath, 0, -1)
		require.Truef(t, errors.Is(err, ErrUnsupportedOffsetLimit), "actual err - %v", err)
	})
	t.Run(ErrOffsetExceedsFileSize.Error(), func(t *testing.T) {
		_, fromPath := createTempFile(t)
		_, toPath := createTempFile(t)
		_, err := Copy(fromPath, toPath, 5, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})
	t.Run(ErrToFileDirNotExists.Error(), func(t *testing.T) {
		_, fromPath := createTempFile(t)
		_, err := Copy(fromPath, "", 0, 0)
		require.Truef(t, errors.Is(err, ErrToFileDirNotExists), "actual err - %v", err)
	})
	t.Run(ErrToFileOpen.Error(), func(t *testing.T) {
		_, fromPath := createTempFile(t)
		f, err := os.OpenFile(".tmp", os.O_RDWR|os.O_CREATE, os.ModeExclusive)
		require.NoError(t, err)
		defer removeAndCloseTempFile(f)
		_, err = Copy(fromPath, f.Name(), 0, 0)
		require.Truef(t, errors.Is(err, ErrToFileOpen), "actual err - %v", err)
	})
	t.Run("regular file", func(t *testing.T) {
		fromFile, fromPath := createTempFile(t)
		toFile, toPath := createTempFile(t)
		subTestTearDown := func() {
			_, err := toFile.Seek(0, io.SeekStart)
			require.NoError(t, err)
		}
		fromString := "Test_String"
		bytesToWrite := int64(len(fromString))
		_, err := io.WriteString(fromFile, fromString)
		require.NoError(t, err)
		offset := bytesToWrite / 2
		limit := bytesToWrite - offset
		t.Run("no offset-limit", func(t *testing.T) {
			t.Cleanup(subTestTearDown)
			bytesWritten, err := Copy(fromPath, toPath, 0, 0)
			require.NoError(t, err)
			require.Equal(t, bytesToWrite, bytesWritten, "not all bytes written")
			toStringBuf, err := io.ReadAll(toFile)
			require.NoError(t, err)
			toString := string(toStringBuf)
			require.Equal(t, fromString, toString, "content not identical")
		})
		t.Run("with offset only", func(t *testing.T) {
			t.Cleanup(subTestTearDown)
			bytesWritten, err := Copy(fromPath, toPath, offset, 0)
			require.NoError(t, err)
			require.Equal(t, bytesToWrite-offset, bytesWritten, "not all bytes written")
			toStringBuf, err := io.ReadAll(toFile)
			require.NoError(t, err)
			toString := string(toStringBuf)
			require.Equal(t, fromString[offset:], toString, "content not identical")
		})
		t.Run("with limit only", func(t *testing.T) {
			t.Cleanup(subTestTearDown)
			bytesWritten, err := Copy(fromPath, toPath, 0, limit)
			require.NoError(t, err)
			require.Equal(t, limit, bytesWritten, "not all bytes written")
			toStringBuf, err := io.ReadAll(toFile)
			require.NoError(t, err)
			toString := string(toStringBuf)
			require.Equal(t, fromString[0:limit], toString, "content not identical")
		})
		t.Run("with offset-limit", func(t *testing.T) {
			t.Cleanup(subTestTearDown)
			bytesWritten, err := Copy(fromPath, toPath, offset, limit-1)
			require.NoError(t, err)
			require.Equal(t, limit-1, bytesWritten, "not all bytes written")
			toStringBuf, err := io.ReadAll(toFile)
			require.NoError(t, err)
			toString := string(toStringBuf)
			require.Equal(t, fromString[offset:offset+limit-1], toString, "content not identical")
		})
		t.Run("with offset-limit and overhead", func(t *testing.T) {
			t.Cleanup(subTestTearDown)
			bytesWritten, err := Copy(fromPath, toPath, offset*2, limit)
			require.NoError(t, err)
			require.Equal(t, bytesToWrite-offset*2, bytesWritten, "not all bytes written")
			toStringBuf, err := io.ReadAll(toFile)
			require.NoError(t, err)
			toString := string(toStringBuf)
			require.Equal(t, fromString[offset*2:limit+offset], toString, "content not identical")
		})
	})
	t.Run("irregular file", func(t *testing.T) {
		if isUnixLike() {
			_, toPath := createTempFile(t)
			t.Run("without limit", func(t *testing.T) {
				_, err := Copy("/dev/null", toPath, 0, 0)
				require.Truef(t, errors.Is(err, ErrFromFileUnsupported), "actual err - %v", err)
			})
			t.Run("with limit", func(t *testing.T) {
				limit := int64(5)
				bytesWritten, err := Copy("/dev/urandom", toPath, 0, limit)
				require.NoError(t, err)
				require.Equal(t, limit, bytesWritten, "not all bytes written")
			})
		}
	})
}
