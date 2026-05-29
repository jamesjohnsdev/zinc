// Package git wraps the git CLI to provide the operations zinc's UI needs:
// status, staging, diffs, commits, branches, log, and stash.
package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes git commands rooted at a working directory.
type Runner struct {
	// Dir is the directory git commands are run from.
	Dir string
}

// New returns a Runner rooted at dir.
func New(dir string) *Runner {
	return &Runner{Dir: dir}
}

// run executes git with the given arguments and returns stdout. On failure
// it returns an error containing the command and git's stderr output.
func (r *Runner) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = r.Dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}

	return stdout.String(), nil
}

// RootDir returns the top-level directory of the git working tree
// containing dir, or an error if dir is not inside one.
func RootDir(dir string) (string, error) {
	out, err := (&Runner{Dir: dir}).run(context.Background(), "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
