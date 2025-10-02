package upd8

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// CommandResult contains stdout, stderr, and exit status from a command execution.
type CommandResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Error    error
}

// CommandRunner executes commands and returns a CommandResult.
type CommandRunner interface {
	Run(ctx context.Context, cmd string, args ...string) CommandResult
}

// ExecRunner executes commands using the local OS shell.
type ExecRunner struct {
	Timeout time.Duration
}

// Run executes a command with a context-aware timeout and captures stdout/stderr.
func (r ExecRunner) Run(ctx context.Context, cmd string, args ...string) CommandResult {
	if r.Timeout <= 0 {
		r.Timeout = 45 * time.Second
	}

	nCtx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	command := exec.CommandContext(nCtx, cmd, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	res := CommandResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), Error: err}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				res.ExitCode = waitStatus.ExitStatus()
			}
		} else if errors.Is(err, context.DeadlineExceeded) {
			res.ExitCode = -1
			res.Error = fmt.Errorf("command timed out: %w", err)
		}
	} else {
		res.ExitCode = 0
	}
	return res
}

// trimStdout normalises stdout/stderr for logging/diagnostics.
func trimStdout(data []byte) string {
	return strings.TrimSpace(string(data))
}

// CombineOutput returns stdout+stderr trimmed, useful for quick summaries.
func (r CommandResult) CombineOutput() string {
	return strings.TrimSpace(trimStdout(r.Stdout) + "\n" + trimStdout(r.Stderr))
}

// HasOutput reports if the command produced any stdout payload.
func (r CommandResult) HasOutput() bool {
	return len(bytes.TrimSpace(r.Stdout)) > 0
}
