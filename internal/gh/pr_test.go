package gh

import (
	"encoding/json"
	"testing"
	"time"
)

// sample captured from a real `gh pr list --json ...` invocation.
const samplePRList = `[{"author":{"id":"MDQ6VXNlcjQ3Mzk0MjAw","is_bot":false,"login":"BagToad","name":"Kynan Ware"},"baseRefName":"bagtoad/attach-issue-commands","headRefName":"bagtoad/attach-review-feedback","isDraft":true,"labels":[],"number":14200,"state":"OPEN","title":"` + "`--attach` stack: review feedback round" + `","updatedAt":"2026-08-19T22:54:44Z"}]`

func TestPullRequestUnmarshal(t *testing.T) {
	var prs []PullRequest
	if err := json.Unmarshal([]byte(samplePRList), &prs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(prs) != 1 {
		t.Fatalf("got %d PRs, want 1", len(prs))
	}

	pr := prs[0]
	if pr.Number != 14200 {
		t.Errorf("Number = %d, want 14200", pr.Number)
	}
	if pr.Author.Login != "BagToad" {
		t.Errorf("Author.Login = %q, want %q", pr.Author.Login, "BagToad")
	}
	if pr.State != "OPEN" {
		t.Errorf("State = %q, want OPEN", pr.State)
	}
	if !pr.IsDraft {
		t.Error("IsDraft = false, want true")
	}
	if pr.HeadRefName != "bagtoad/attach-review-feedback" {
		t.Errorf("HeadRefName = %q", pr.HeadRefName)
	}
	if pr.BaseRefName != "bagtoad/attach-issue-commands" {
		t.Errorf("BaseRefName = %q", pr.BaseRefName)
	}
	want := time.Date(2026, 8, 19, 22, 54, 44, 0, time.UTC)
	if !pr.UpdatedAt.Equal(want) {
		t.Errorf("UpdatedAt = %v, want %v", pr.UpdatedAt, want)
	}
	if len(pr.Labels) != 0 {
		t.Errorf("Labels = %v, want empty", pr.Labels)
	}
}

func TestPullRequestUnmarshalWithLabels(t *testing.T) {
	in := `{"number":1,"title":"t","author":{"login":"me"},"state":"OPEN","isDraft":false,"labels":[{"name":"bug"},{"name":"needs-triage"}]}`

	var pr PullRequest
	if err := json.Unmarshal([]byte(in), &pr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(pr.Labels) != 2 || pr.Labels[0].Name != "bug" || pr.Labels[1].Name != "needs-triage" {
		t.Errorf("Labels = %+v", pr.Labels)
	}
}
