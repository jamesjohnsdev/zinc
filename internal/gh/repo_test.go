package gh

import (
	"encoding/json"
	"testing"
)

const sampleRepoView = `{"defaultBranchRef":{"name":"trunk"},"description":"GitHub's official command line tool","isPrivate":false,"name":"cli","nameWithOwner":"cli/cli","owner":{"login":"cli"},"stargazerCount":45895,"url":"https://github.com/cli/cli"}`

func TestRepoUnmarshal(t *testing.T) {
	var repo Repo
	if err := json.Unmarshal([]byte(sampleRepoView), &repo); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if repo.Name != "cli" {
		t.Errorf("Name = %q, want cli", repo.Name)
	}
	if repo.NameWithOwner != "cli/cli" {
		t.Errorf("NameWithOwner = %q", repo.NameWithOwner)
	}
	if repo.Owner.Login != "cli" {
		t.Errorf("Owner.Login = %q", repo.Owner.Login)
	}
	if repo.DefaultBranch.Name != "trunk" {
		t.Errorf("DefaultBranch.Name = %q, want trunk", repo.DefaultBranch.Name)
	}
	if repo.IsPrivate {
		t.Error("IsPrivate = true, want false")
	}
	if repo.StargazerCount != 45895 {
		t.Errorf("StargazerCount = %d", repo.StargazerCount)
	}
}
