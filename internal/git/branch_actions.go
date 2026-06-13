package git

import (
	"context"
	"strings"
)

// Checkout switches the working tree to the given local branch.
func (r *Runner) Checkout(ctx context.Context, name string) error {
	_, err := r.run(ctx, "checkout", name)
	return err
}

// CheckoutRemote creates and checks out a local branch tracking the given
// remote-tracking branch (e.g. "origin/feature" creates local "feature").
func (r *Runner) CheckoutRemote(ctx context.Context, remoteBranch string) error {
	local := remoteBranch
	if i := strings.Index(remoteBranch, "/"); i >= 0 {
		local = remoteBranch[i+1:]
	}

	_, err := r.run(ctx, "checkout", "-b", local, "--track", remoteBranch)
	return err
}

// CreateBranch creates and checks out a new branch from the current HEAD.
func (r *Runner) CreateBranch(ctx context.Context, name string) error {
	_, err := r.run(ctx, "checkout", "-b", name)
	return err
}

// DeleteBranch deletes a local branch. Unless force is set, git refuses to
// delete a branch that isn't fully merged.
func (r *Runner) DeleteBranch(ctx context.Context, name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}

	_, err := r.run(ctx, "branch", flag, name)
	return err
}
