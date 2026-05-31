package git

import (
	"context"
	"strings"
)

// FileStatus describes the state of a single path in the working tree and
// index, as reported by `git status --porcelain=v2`.
type FileStatus struct {
	// Path is the file's current path, relative to the repository root.
	Path string
	// OrigPath is set for renames/copies to the path's previous location.
	OrigPath string
	// Staged is the index status code ('.' if unchanged in the index).
	Staged byte
	// Unstaged is the worktree status code ('.' if unchanged in the worktree).
	Unstaged byte
	// Untracked is true for paths git does not yet track.
	Untracked bool
	// Unmerged is true for paths with unresolved merge conflicts.
	Unmerged bool
}

// IsRename reports whether the entry represents a rename or copy.
func (f FileStatus) IsRename() bool {
	return f.OrigPath != ""
}

// Status returns the working tree and index status of every changed,
// untracked, or unmerged path in the repository.
func (r *Runner) Status(ctx context.Context) ([]FileStatus, error) {
	out, err := r.run(ctx, "status", "--porcelain=v2", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	return parsePorcelainV2(out), nil
}

func parsePorcelainV2(out string) []FileStatus {
	tokens := strings.Split(out, "\x00")

	var files []FileStatus
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if tok == "" {
			continue
		}

		switch tok[0] {
		case '1': // ordinary changed entry
			fields := strings.SplitN(tok, " ", 9)
			if len(fields) < 9 {
				continue
			}
			xy := fields[1]
			files = append(files, FileStatus{
				Path:     fields[8],
				Staged:   xy[0],
				Unstaged: xy[1],
			})

		case '2': // renamed or copied entry
			fields := strings.SplitN(tok, " ", 10)
			if len(fields) < 10 {
				continue
			}
			xy := fields[1]
			path := fields[9]

			var orig string
			if i+1 < len(tokens) {
				i++
				orig = tokens[i]
			}

			files = append(files, FileStatus{
				Path:     path,
				OrigPath: orig,
				Staged:   xy[0],
				Unstaged: xy[1],
			})

		case 'u': // unmerged entry
			fields := strings.SplitN(tok, " ", 11)
			if len(fields) < 11 {
				continue
			}
			xy := fields[1]
			files = append(files, FileStatus{
				Path:     fields[10],
				Staged:   xy[0],
				Unstaged: xy[1],
				Unmerged: true,
			})

		case '?': // untracked
			files = append(files, FileStatus{
				Path:      strings.TrimPrefix(tok, "? "),
				Untracked: true,
				Staged:    '?',
				Unstaged:  '?',
			})

		case '!': // ignored; not requested, but ignore defensively
			continue
		}
	}

	return files
}
