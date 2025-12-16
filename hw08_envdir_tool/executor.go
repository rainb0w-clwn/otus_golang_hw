package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return ExitCode
	}
	err := SetEnv(&env)
	if err != nil {
		return ExitCode
	}
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	c := exec.Command(cmd[0], args...) //nolint:gosec
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return ExitCode
	}
	return
}

func SetEnv(env *Environment) error {
	for k, v := range *env {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return err
			}
		} else {
			if err := os.Setenv(k, v.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
