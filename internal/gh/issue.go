package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Issue is a GitHub issue as shown in list and detail views.
type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Author    Author    `json:"author"`
	State     string    `json:"state"` // OPEN, CLOSED
	Labels    []Label   `json:"labels"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	URL       string    `json:"url"`
}

const (
	issueListFields   = "number,title,author,state,labels,updatedAt"
	issueDetailFields = issueListFields + ",body,url"
)

// Issues lists open issues for repo (or the repo resolved from the working
// directory if repo is empty).
func (r *Runner) Issues(ctx context.Context, repo string) ([]Issue, error) {
	args := append([]string{"issue", "list", "--json", issueListFields, "--limit", "200"}, repoArgs(repo)...)

	out, err := r.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := json.Unmarshal([]byte(out), &issues); err != nil {
		return nil, fmt.Errorf("parse issue list: %w", err)
	}

	return issues, nil
}

// ViewIssue fetches full detail, including body, for a single issue.
func (r *Runner) ViewIssue(ctx context.Context, repo string, number int) (Issue, error) {
	args := append([]string{"issue", "view", strconv.Itoa(number), "--json", issueDetailFields}, repoArgs(repo)...)

	out, err := r.run(ctx, args...)
	if err != nil {
		return Issue{}, err
	}

	var issue Issue
	if err := json.Unmarshal([]byte(out), &issue); err != nil {
		return Issue{}, fmt.Errorf("parse issue view: %w", err)
	}

	return issue, nil
}

// IssueClose closes an issue.
func (r *Runner) IssueClose(ctx context.Context, repo string, number int) error {
	args := append([]string{"issue", "close", strconv.Itoa(number)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// IssueReopen reopens a closed issue.
func (r *Runner) IssueReopen(ctx context.Context, repo string, number int) error {
	args := append([]string{"issue", "reopen", strconv.Itoa(number)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// IssueComment adds a comment to an issue.
func (r *Runner) IssueComment(ctx context.Context, repo string, number int, body string) error {
	args := append([]string{"issue", "comment", strconv.Itoa(number), "--body", body}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// IssueWeb opens an issue in the default web browser.
func (r *Runner) IssueWeb(ctx context.Context, repo string, number int) error {
	args := append([]string{"issue", "view", strconv.Itoa(number), "--web"}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}
