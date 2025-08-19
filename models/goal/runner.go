package goal

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner struct {
	Binary          string        // default "goal"
	DataDirs        []string      // multiple -d allowed
	KmdDir          string        // -k
	Wallet          string        // -w
	Timeout         time.Duration // default 20s
	AllowEmptyWallet bool
}

func NewRunner() *Runner {
	return &Runner{
		Binary:  "goal",
		Timeout: 20 * time.Second,
	}
}

// baseFlags composes [-d ...] [-k ...] [-w ...] for any goal command.
func (r *Runner) baseFlags() []string {
	var flags []string
	for _, d := range r.DataDirs {
		if strings.TrimSpace(d) != "" {
			flags = append(flags, "-d", d)
		}
	}
	if strings.TrimSpace(r.KmdDir) != "" {
		flags = append(flags, "-k", r.KmdDir)
	}
	if strings.TrimSpace(r.Wallet) != "" {
		flags = append(flags, "-w", r.Wallet)
	}
	return flags
}

type RunResult struct {
	Stdout string
	Stderr string
	Err    error
}

func (r *Runner) Run(ctx context.Context, argv []string) RunResult {
	if r.Timeout <= 0 {
		r.Timeout = 20 * time.Second
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, r.Binary, append(r.baseFlags(), argv...)...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	cmd.Env = os.Environ()

	err := cmd.Run()
	return RunResult{Stdout: out.String(), Stderr: errb.String(), Err: err}
}

func (r *Runner) CheckBinary() error {
	_, err := exec.LookPath(r.Binary)
	if err != nil {
		return errors.New("`goal` binary not found in PATH")
	}
	return nil
}
