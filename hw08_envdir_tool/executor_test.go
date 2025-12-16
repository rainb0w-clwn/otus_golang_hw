package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var val = "abc"

func TestRunCmd(t *testing.T) {
	t.Run("no cmd error", func(t *testing.T) {
		res := RunCmd([]string{}, Environment{})
		require.Equal(t, ExitCode, res)
	})

	t.Run("env error", func(t *testing.T) {
		res := RunCmd([]string{val}, Environment{
			"": EnvValue{},
		})
		require.Equal(t, ExitCode, res)
	})

	t.Run("cmd syserror", func(t *testing.T) {
		res := RunCmd([]string{""}, Environment{
			val: EnvValue{},
		})
		require.Equal(t, ExitCode, res)
	})

	t.Run("cmd with args", func(t *testing.T) {
		command := "echo"
		args := []string{"-n", "arg1", "arg2"}
		cmd := append([]string{command}, args...)
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		res := RunCmd(cmd, Environment{
			val: EnvValue{},
		})
		require.Equal(t, 0, res)
		err := w.Close()
		require.NoError(t, err)
		out, _ := io.ReadAll(r)
		os.Stdout = old
		require.Equal(t, args[1]+" "+args[2], string(out))
	})

	t.Run("provide cmd exitCode", func(t *testing.T) {
		command := "/bin/bash"
		args := []string{"a"}
		cmd := append([]string{command}, args...)
		old := os.Stderr

		// Get expected stdErr bytes
		r, w, _ := os.Pipe()
		c := exec.Command(cmd[0], args...)
		c.Stderr = w
		var exitError *exec.ExitError
		err := c.Run()
		require.Error(t, err)
		errors.As(err, &exitError)
		outE := make([]byte, 5)
		_, err = io.ReadFull(r, outE)
		require.NoError(t, err)
		err = w.Close()
		require.NoError(t, err)

		// Set new stdErr
		r, w, _ = os.Pipe()
		os.Stderr = w
		res := RunCmd(cmd, Environment{
			val: EnvValue{},
		})
		err = w.Close()
		require.NoError(t, err)
		require.Equal(t, exitError.ExitCode(), res)

		// Get actual stdErr bytes
		outA := make([]byte, 5)
		_, err = io.ReadFull(r, outA)
		os.Stderr = old
		require.NoError(t, err)
		require.Equal(t, outE, outA)
	})

	t.Run("provide stdIn to cmd", func(t *testing.T) {
		command := "/bin/bash"
		expected := val + "\n"
		c := make(chan interface{})
		old := os.Stdout
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut
		oldStdin := os.Stdin
		rIn, wIn, _ := os.Pipe()
		os.Stdin = rIn
		_, err := wIn.WriteString("echo " + expected)
		require.NoError(t, err)
		go func() {
			time.Sleep(500 * time.Millisecond)
			_, err = wIn.WriteString("exit\n")
			require.NoError(t, err)
			c <- nil
		}()
		_ = RunCmd([]string{command}, Environment{
			val: EnvValue{},
		})
		<-c
		err = wOut.Close()
		require.NoError(t, err)
		err = rIn.Close()
		require.NoError(t, err)
		s, err := io.ReadAll(rOut)
		require.NoError(t, err)
		os.Stdout = old
		os.Stdin = oldStdin
		require.Equal(t, expected, string(s))
	})
}

func TestSetEnv(t *testing.T) {
	t.Cleanup(func() {
		os.Clearenv()
	})
	t.Run("remove", func(t *testing.T) {
		key := "RANDOM_ENV_KEY"
		err := os.Setenv(key, val)
		require.NoError(t, err)
		_, exist := os.LookupEnv(key)
		require.True(t, exist)
		err = SetEnv(&Environment{
			key: EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		})
		require.NoError(t, err)
		_, exist = os.LookupEnv(key)
		require.False(t, exist)
	})
	t.Run("set", func(t *testing.T) {
		key := "RANDOM_ENV_KEY"
		err := os.Setenv(key, val)
		require.NoError(t, err)
		v := os.Getenv(key)
		require.Equal(t, val, v)
		err = SetEnv(&Environment{
			key: EnvValue{
				Value:      "",
				NeedRemove: false,
			},
		})
		require.NoError(t, err)
		v = os.Getenv(key)
		require.Equal(t, "", v)
	})
	t.Run("combined", func(t *testing.T) {
		key1 := "RANDOM_ENV_KEY_1"
		key2 := "RANDOM_ENV_KEY_2"
		err := os.Setenv(key1, val)
		require.NoError(t, err)
		v := os.Getenv(key1)
		require.Equal(t, val, v)
		err = os.Setenv(key2, val)
		require.NoError(t, err)
		v = os.Getenv(key2)
		require.Equal(t, val, v)
		err = SetEnv(&Environment{
			key1: EnvValue{
				Value:      val,
				NeedRemove: true,
			},
			key2: EnvValue{
				Value:      val,
				NeedRemove: false,
			},
		})
		require.NoError(t, err)
		_, exist := os.LookupEnv(key1)
		require.False(t, exist)
		v = os.Getenv(key2)
		require.Equal(t, val, v)
	})
}
