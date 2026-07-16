package gh

import (
	"encoding/json"
	"testing"
)

const sampleIssueList = `[{"author":{"id":"MDQ6VXNlcjU0NTEzMDY2","is_bot":false,"login":"LeonardoLGDS","name":"leosiedler"},"labels":[{"id":"LA_kwDODKw3uc7QD3p7","name":"needs-triage","description":"needs to be reviewed","color":"D6393F"}],"number":14199,"state":"OPEN","title":"gh api: warn when -f value starts with @ and names an existing file","updatedAt":"2026-08-19T15:22:45Z"}]`

func TestIssueUnmarshal(t *testing.T) {
	var issues []Issue
	if err := json.Unmarshal([]byte(sampleIssueList), &issues); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("got %d issues, want 1", len(issues))
	}

	issue := issues[0]
	if issue.Number != 14199 {
		t.Errorf("Number = %d, want 14199", issue.Number)
	}
	if issue.Author.Login != "LeonardoLGDS" {
		t.Errorf("Author.Login = %q", issue.Author.Login)
	}
	if issue.State != "OPEN" {
		t.Errorf("State = %q, want OPEN", issue.State)
	}
	if len(issue.Labels) != 1 || issue.Labels[0].Name != "needs-triage" {
		t.Errorf("Labels = %+v", issue.Labels)
	}
}
