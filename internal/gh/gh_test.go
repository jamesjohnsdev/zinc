package gh

import "testing"

func TestRepoArgs(t *testing.T) {
	if got := repoArgs(""); got != nil {
		t.Errorf("repoArgs(%q) = %v, want nil", "", got)
	}

	got := repoArgs("owner/name")
	want := []string{"--repo", "owner/name"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("repoArgs(%q) = %v, want %v", "owner/name", got, want)
	}
}
