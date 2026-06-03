package git

import "context"

// Stage adds path's current worktree contents to the index.
func (r *Runner) Stage(ctx context.Context, path string) error {
	_, err := r.run(ctx, "add", "--", path)
	return err
}

// StageAll stages every change in the working tree, including untracked
// files.
func (r *Runner) StageAll(ctx context.Context) error {
	_, err := r.run(ctx, "add", "-A")
	return err
}

// Unstage removes path from the index without touching the working tree.
func (r *Runner) Unstage(ctx context.Context, path string) error {
	_, err := r.run(ctx, "restore", "--staged", "--", path)
	return err
}

// UnstageAll unstages every currently staged change.
func (r *Runner) UnstageAll(ctx context.Context) error {
	_, err := r.run(ctx, "restore", "--staged", ".")
	return err
}

// Discard reverts path's working tree changes. Untracked paths are removed
// via clean; tracked paths are restored from the index.
func (r *Runner) Discard(ctx context.Context, path string, untracked bool) error {
	if untracked {
		_, err := r.run(ctx, "clean", "-f", "--", path)
		return err
	}
	_, err := r.run(ctx, "checkout", "--", path)
	return err
}
