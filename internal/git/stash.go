package git

import (
	"context"
	"fmt"
	"strings"
)

// Stash describes a single entry in the stash list.
type Stash struct {
	Index   int
	Message string
}

// StashList returns the stash entries, most recent (stash@{0}) first.
func (r *Runner) StashList(ctx context.Context) ([]Stash, error) {
	out, err := r.run(ctx, "stash", "list", "--format=%gs")
	if err != nil {
		return nil, err
	}

	var stashes []Stash
	for i, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if line == "" {
			continue
		}
		stashes = append(stashes, Stash{Index: i, Message: line})
	}

	return stashes, nil
}

// StashPush stashes the current working tree changes. If includeUntracked
// is set, untracked files are stashed too.
func (r *Runner) StashPush(ctx context.Context, message string, includeUntracked bool) error {
	args := []string{"stash", "push"}
	if includeUntracked {
		args = append(args, "-u")
	}
	if message != "" {
		args = append(args, "-m", message)
	}

	_, err := r.run(ctx, args...)
	return err
}

// StashPop applies and removes the stash entry at index.
func (r *Runner) StashPop(ctx context.Context, index int) error {
	_, err := r.run(ctx, "stash", "pop", fmt.Sprintf("stash@{%d}", index))
	return err
}

// StashDrop removes the stash entry at index without applying it.
func (r *Runner) StashDrop(ctx context.Context, index int) error {
	_, err := r.run(ctx, "stash", "drop", fmt.Sprintf("stash@{%d}", index))
	return err
}
