package pdf

import (
	"bytes"
	"os"
	"os/exec"
	"time"
)

// RunCommand performs run of desired command with timeout option
func RunCommand(dir, execFile string, maxRuntime time.Duration, environment []string, args ...string) (stdout, stderr string, err error) {

	cmd := exec.Command(execFile, args...)
	var (
		stdoutB bytes.Buffer
		stderrB bytes.Buffer
	)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, environment...)

	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB

	if dir != "" {
		cmd.Dir = dir
	}

	done := make(chan error)
	go func() {
		cmd.Start()
		done <- cmd.Wait()
	}()
	defer func() {
		cmd.Process.Kill()
	}()

	select {
	case <-time.After(maxRuntime):
		cmd.Process.Kill()
		return

	case err = <-done:
		stdout = stdoutB.String()
		stderr = stderrB.String()
	}

	return
}
