package git

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Commit is a single entry in the commit log.
type Commit struct {
	Hash    string
	Short   string
	Author  string
	Date    time.Time
	Subject string
}

// Log returns the commit history reachable from HEAD, most recent first.
// A non-positive limit requests the full history.
func (r *Runner) Log(ctx context.Context, limit int) ([]Commit, error) {
	format := "%H%x09%h%x09%an%x09%aI%x09%s"
	args := []string{"log", "--format=" + format}
	if limit > 0 {
		args = append(args, fmt.Sprintf("-%d", limit))
	}

	out, err := r.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.SplitN(line, "\t", 5)
		if len(fields) < 5 {
			continue
		}

		date, _ := time.Parse(time.RFC3339, fields[3])

		commits = append(commits, Commit{
			Hash:    fields[0],
			Short:   fields[1],
			Author:  fields[2],
			Date:    date,
			Subject: fields[4],
		})
	}

	return commits, nil
}
