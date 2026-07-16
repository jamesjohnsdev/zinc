// Package gh wraps the gh CLI to provide the GitHub operations zinc's UI
// needs: repo info, pull requests, issues, and Actions workflow runs.
package gh

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes gh commands rooted at a working directory.
type Runner struct {
	// Dir is the directory gh commands are run from. gh resolves the
	// target repository from the git remote in this directory unless a
	// repo is explicitly passed to a method.
	Dir string
}

// New returns a Runner rooted at dir.
func New(dir string) *Runner {
	return &Runner{Dir: dir}
}

// run executes gh with the given arguments and returns stdout. On failure
// it returns an error containing the command and gh's stderr output.
func (r *Runner) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Dir = r.Dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("gh %s: %s", strings.Join(args, " "), msg)
	}

	return stdout.String(), nil
}

// repoArgs returns the --repo flag for repo, or nil to let gh resolve the
// repository from the working directory's git remote. repo is an
// "[HOST/]OWNER/REPO" string.
func repoArgs(repo string) []string {
	if repo == "" {
		return nil
	}
	return []string{"--repo", repo}
}

// Author is a GitHub user, as returned in PR/issue "author" JSON fields.
type Author struct {
	Login string `json:"login"`
}

// Label is an issue/PR label.
type Label struct {
	Name string `json:"name"`
}
