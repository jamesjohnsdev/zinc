package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// WorkflowRun is a single run of a GitHub Actions workflow.
type WorkflowRun struct {
	DatabaseID   int64     `json:"databaseId"`
	Name         string    `json:"name"`
	WorkflowName string    `json:"workflowName"`
	Status       string    `json:"status"`     // queued, in_progress, completed
	Conclusion   string    `json:"conclusion"` // success, failure, cancelled, ... ("" while not completed)
	HeadBranch   string    `json:"headBranch"`
	Event        string    `json:"event"`
	CreatedAt    time.Time `json:"createdAt"`
	URL          string    `json:"url"`
}

const runListFields = "databaseId,name,workflowName,status,conclusion,headBranch,event,createdAt,url"

// WorkflowRuns lists recent Actions workflow runs for repo (or the repo
// resolved from the working directory if repo is empty).
func (r *Runner) WorkflowRuns(ctx context.Context, repo string) ([]WorkflowRun, error) {
	args := append([]string{"run", "list", "--json", runListFields, "--limit", "50"}, repoArgs(repo)...)

	out, err := r.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var runs []WorkflowRun
	if err := json.Unmarshal([]byte(out), &runs); err != nil {
		return nil, fmt.Errorf("parse run list: %w", err)
	}

	return runs, nil
}

// RunRerun re-runs a workflow run.
func (r *Runner) RunRerun(ctx context.Context, repo string, id int64) error {
	args := append([]string{"run", "rerun", strconv.FormatInt(id, 10)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}

// RunCancel cancels an in-progress workflow run.
func (r *Runner) RunCancel(ctx context.Context, repo string, id int64) error {
	args := append([]string{"run", "cancel", strconv.FormatInt(id, 10)}, repoArgs(repo)...)
	_, err := r.run(ctx, args...)
	return err
}
