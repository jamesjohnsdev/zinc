package gh

import (
	"encoding/json"
	"testing"
)

const sampleRunList = `[{"conclusion":"failure","createdAt":"2026-08-20T04:53:40Z","databaseId":32333541211,"event":"schedule","headBranch":"trunk","name":"Dependabot PR Triage (skills-driven)","status":"completed","workflowName":"Dependabot PR Triage (skills-driven)","url":"https://github.com/cli/cli/actions/runs/32333541211"}]`

func TestWorkflowRunUnmarshal(t *testing.T) {
	var runs []WorkflowRun
	if err := json.Unmarshal([]byte(sampleRunList), &runs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("got %d runs, want 1", len(runs))
	}

	run := runs[0]
	if run.DatabaseID != 32333541211 {
		t.Errorf("DatabaseID = %d", run.DatabaseID)
	}
	if run.Status != "completed" {
		t.Errorf("Status = %q, want completed", run.Status)
	}
	if run.Conclusion != "failure" {
		t.Errorf("Conclusion = %q, want failure", run.Conclusion)
	}
	if run.HeadBranch != "trunk" {
		t.Errorf("HeadBranch = %q, want trunk", run.HeadBranch)
	}
}
