package git

import "context"

// Commit records the currently staged changes with the given message.
func (r *Runner) Commit(ctx context.Context, message string) error {
	_, err := r.run(ctx, "commit", "-m", message)
	return err
}
