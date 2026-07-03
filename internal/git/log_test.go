package git

import (
	"testing"
	"time"
)

func TestParseLog(t *testing.T) {
	in := "abc123\tabc\tJane Doe\t2026-06-01T10:00:00+00:00\tfirst commit\n" +
		"def456\tdef\tJohn Roe\t2026-06-02T11:30:00+00:00\tfix: second commit\n"

	got := parseLog(in)
	if len(got) != 2 {
		t.Fatalf("got %d commits, want 2: %+v", len(got), got)
	}

	want0 := Commit{
		Hash:    "abc123",
		Short:   "abc",
		Author:  "Jane Doe",
		Date:    time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
		Subject: "first commit",
	}
	if !got[0].Date.Equal(want0.Date) {
		t.Errorf("date = %v, want %v", got[0].Date, want0.Date)
	}
	got[0].Date = want0.Date
	if got[0] != want0 {
		t.Errorf("commit 0 = %+v, want %+v", got[0], want0)
	}

	if got[1].Subject != "fix: second commit" {
		t.Errorf("commit 1 subject = %q, want %q", got[1].Subject, "fix: second commit")
	}
}

func TestParseLogEmpty(t *testing.T) {
	if got := parseLog(""); got != nil {
		t.Errorf("got %+v, want nil", got)
	}
}
