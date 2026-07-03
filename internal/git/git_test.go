package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// newTestRepo creates an empty git repository in a temporary directory and
// returns a Runner rooted at it.
func newTestRepo(t *testing.T) *Runner {
	t.Helper()

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init", "-q")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")

	return New(dir)
}

func TestRunnerStageAndCommit(t *testing.T) {
	r := newTestRepo(t)
	ctx := context.Background()

	path := filepath.Join(r.Dir, "file.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := r.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || !files[0].Untracked {
		t.Fatalf("status = %+v, want one untracked file", files)
	}

	if err := r.Stage(ctx, "file.txt"); err != nil {
		t.Fatal(err)
	}

	files, err = r.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Staged != 'A' {
		t.Fatalf("status after stage = %+v, want a staged add", files)
	}

	if err := r.Commit(ctx, "add file"); err != nil {
		t.Fatal(err)
	}

	files, err = r.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("status after commit = %+v, want a clean tree", files)
	}

	commits, err := r.Log(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 || commits[0].Subject != "add file" {
		t.Fatalf("log = %+v, want a single 'add file' commit", commits)
	}
}
