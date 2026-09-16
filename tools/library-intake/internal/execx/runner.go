package execx

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	derrors "github.com/xynova/library-intake/internal/errors"
)

// Runner executes external processes with context cancellation.
type Runner struct {
	LookPath func(file string) (string, error)
	Command  func(ctx context.Context, name string, args ...string) *exec.Cmd
}

// NewRunner returns a Runner that uses the OS look-path and command factories.
func NewRunner() *Runner {
	return &Runner{
		LookPath: exec.LookPath,
		Command:  exec.CommandContext,
	}
}

// Result is captured process output.
type Result struct {
	Stdout string
	Stderr string
}

// Run executes name with args, returning wrapped errors that include stderr.
func (r *Runner) Run(ctx context.Context, name string, args ...string) (*Result, error) {
	if r == nil {
		return nil, derrors.New(derrors.CodeFailed, "execx.Run", "runner is nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, derrors.Wrap(err, derrors.CodeFailed, "execx.Run", "context done")
	}
	path, err := r.LookPath(name)
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeNotFound, "execx.Run", "binary not found").
			With("binary", name)
	}
	cmdFactory := r.Command
	if cmdFactory == nil {
		cmdFactory = exec.CommandContext
	}
	cmd := cmdFactory(ctx, path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &Result{Stdout: stdout.String(), Stderr: stderr.String()},
			derrors.Wrap(err, derrors.CodeFailed, "execx.Run", "command failed").
				With("binary", name).
				With("args", strings.Join(args, " ")).
				With("stderr", trimErr(stderr.String()))
	}
	return &Result{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

// RunDir is like Run but sets the process working directory.
func (r *Runner) RunDir(ctx context.Context, dir, name string, args ...string) (*Result, error) {
	if r == nil {
		return nil, derrors.New(derrors.CodeFailed, "execx.RunDir", "runner is nil")
	}
	path, err := r.LookPath(name)
	if err != nil {
		return nil, derrors.Wrap(err, derrors.CodeNotFound, "execx.RunDir", "binary not found").
			With("binary", name)
	}
	cmdFactory := r.Command
	if cmdFactory == nil {
		cmdFactory = exec.CommandContext
	}
	cmd := cmdFactory(ctx, path, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)
	if err := cmd.Run(); err != nil {
		return &Result{Stdout: stdout.String(), Stderr: stderr.String()},
			derrors.Wrap(err, derrors.CodeFailed, "execx.RunDir", "command failed").
				With("binary", name).
				With("dir", dir).
				With("args", strings.Join(args, " ")).
				With("stderr", trimErr(stderr.String()))
	}
	return &Result{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

func trimErr(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 2000 {
		return s[:2000] + "…"
	}
	return s
}
