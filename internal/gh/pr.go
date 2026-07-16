package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// PullRequest is a pull request as shown in list and detail views.
type PullRequest struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Author      Author    `json:"author"`
	State       string    `json:"state"` // OPEN, CLOSED, MERGED
	IsDraft     bool      `json:"isDraft"`
	HeadRefName string    `json:"headRefName"`
	BaseRefName string    `json:"baseRefName"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Labels      []Label   `json:"labels"`
	Body        string    `json:"body"`
	URL         string    `json:"url"`
}

const (
	prListFields   = "number,title,author,state,isDraft,headRefName,baseRefName,updatedAt,labels"
	prDetailFields = prListFields + ",body,url"
)

// PullRequests lists open pull requests for repo (or the repo resolved
// from the working directory if repo is empty).
func (r *Runner) PullRequests(ctx context.Context, repo string) ([]PullRequest, error) {
	args := append([]string{"pr", "list", "--json", prListFields, "--limit", "200"}, repoArgs(repo)...)

	out, err := r.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var prs []PullRequest
	if err := json.Unmarshal([]byte(out), &prs); err != nil {
		return nil, fmt.Errorf("parse pr list: %w", err)
	}

	return prs, nil
}

// ViewPR fetches full detail, including body, for a single pull request.
func (r *Runner) ViewPR(ctx context.Context, repo string, number int) (PullRequest, error) {
	args := append([]string{"pr", "view", strconv.Itoa(number), "--json", prDetailFields}, repoArgs(repo)...)

	out, err := r.run(ctx, args...)
	if err != nil {
		return PullRequest{}, err
	}

	var pr PullRequest
	if err := json.Unmarshal([]byte(out), &pr); err != nil {
		return PullRequest{}, fmt.Errorf("parse pr view: %w", err)
	}

	return pr, nil
}

// PRDiff returns the diff for a pull request.
func (r *Runner) PRDiff(ctx context.Context, repo string, number int) (string, error) {
	args := append([]string{"pr", "diff", strconv.Itoa(number)}, repoArgs(repo)...)
	return r.run(ctx, args...)
}

// PRCheckout checks out a pull request's branch locally.
func (r *Runner) PRCheckout(ctx context.Context, repo string, number int) error {
	args := append([]string{"pr", "checkout", strconv.Itoa(number)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// MergeMethod selects how PRMerge merges a pull request.
type MergeMethod string

const (
	MergeCommit MergeMethod = "--merge"
	MergeSquash MergeMethod = "--squash"
	MergeRebase MergeMethod = "--rebase"
)

// PRMerge merges a pull request using the given method.
func (r *Runner) PRMerge(ctx context.Context, repo string, number int, method MergeMethod) error {
	args := append([]string{"pr", "merge", strconv.Itoa(number), string(method)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// PRClose closes a pull request without merging it.
func (r *Runner) PRClose(ctx context.Context, repo string, number int) error {
	args := append([]string{"pr", "close", strconv.Itoa(number)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// PRReopen reopens a closed pull request.
func (r *Runner) PRReopen(ctx context.Context, repo string, number int) error {
	args := append([]string{"pr", "reopen", strconv.Itoa(number)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}
