package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// Diff returns the diff for path. If staged is true, the diff is between
// the index and HEAD (git diff --cached); otherwise it is between the
// working tree and the index.
func (r *Runner) Diff(ctx context.Context, path string, staged bool) (string, error) {
	args := []string{"diff", "--no-color"}
	if staged {
		args = append(args, "--cached")
	}
	args = append(args, "--", path)

	return r.run(ctx, args...)
}

// DiffUntracked returns a diff-style view of an untracked file's full
// contents, computed against /dev/null.
func (r *Runner) DiffUntracked(ctx context.Context, path string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--no-color", "--no-index", "--", "/dev/null", path)
	cmd.Dir = r.Dir

	out, err := cmd.Output()
	if err != nil {
		// git diff --no-index exits 1 when the compared paths differ, which
		// is the expected outcome here, not a failure.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return string(out), nil
		}
		return "", fmt.Errorf("git diff --no-index: %w", err)
	}

	return string(out), nil
}
