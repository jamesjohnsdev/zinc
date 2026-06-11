package git

import (
	"context"
	"fmt"
	"strings"
)

// Branch describes a local or remote-tracking branch.
type Branch struct {
	Name     string
	Remote   bool
	Current  bool
	Upstream string
	Ahead    int
	Behind   int
}

// Branches lists local branches, including their upstream tracking state.
func (r *Runner) Branches(ctx context.Context) ([]Branch, error) {
	format := "%(HEAD)%09%(refname:short)%09%(upstream:short)%09%(upstream:track)"
	out, err := r.run(ctx, "for-each-ref", "--format="+format, "refs/heads/")
	if err != nil {
		return nil, err
	}

	var branches []Branch
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			continue
		}

		b := Branch{
			Name:     fields[1],
			Current:  fields[0] == "*",
			Upstream: fields[2],
		}
		b.Ahead, b.Behind = parseTrack(fields[3])

		branches = append(branches, b)
	}

	return branches, nil
}

// RemoteBranches lists remote-tracking branches (refs/remotes/*), excluding
// symbolic refs such as origin/HEAD.
func (r *Runner) RemoteBranches(ctx context.Context) ([]Branch, error) {
	out, err := r.run(ctx, "for-each-ref", "--format=%(refname:short)", "refs/remotes/")
	if err != nil {
		return nil, err
	}

	var branches []Branch
	for _, name := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if name == "" || strings.HasSuffix(name, "/HEAD") {
			continue
		}
		branches = append(branches, Branch{Name: name, Remote: true})
	}

	return branches, nil
}

// parseTrack parses a for-each-ref %(upstream:track) value such as
// "[ahead 1, behind 2]" into its ahead/behind counts.
func parseTrack(track string) (ahead, behind int) {
	track = strings.Trim(track, "[]")
	if track == "" {
		return 0, 0
	}

	for _, part := range strings.Split(track, ", ") {
		var n int
		if _, err := fmt.Sscanf(part, "ahead %d", &n); err == nil {
			ahead = n
			continue
		}
		if _, err := fmt.Sscanf(part, "behind %d", &n); err == nil {
			behind = n
		}
	}

	return ahead, behind
}
