package git

import "testing"

func TestParsePorcelainV2(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []FileStatus
	}{
		{
			name: "untracked",
			in:   "? new.txt\x00",
			want: []FileStatus{{Path: "new.txt", Untracked: true, Staged: '?', Unstaged: '?'}},
		},
		{
			name: "modified unstaged",
			in:   "1 .M N... 100644 100644 100644 abc abc file.go\x00",
			want: []FileStatus{{Path: "file.go", Staged: '.', Unstaged: 'M'}},
		},
		{
			name: "modified staged",
			in:   "1 M. N... 100644 100644 100644 abc abc file.go\x00",
			want: []FileStatus{{Path: "file.go", Staged: 'M', Unstaged: '.'}},
		},
		{
			name: "renamed",
			in:   "2 R. N... 100644 100644 100644 abc abc R100 new.go\x00old.go\x00",
			want: []FileStatus{{Path: "new.go", OrigPath: "old.go", Staged: 'R', Unstaged: '.'}},
		},
		{
			name: "unmerged",
			in:   "u UU N... 100644 100644 100644 100644 abc abc abc file.go\x00",
			want: []FileStatus{{Path: "file.go", Staged: 'U', Unstaged: 'U', Unmerged: true}},
		},
		{
			name: "multiple entries",
			in:   "1 M. N... 100644 100644 100644 abc abc a.go\x00? b.go\x00",
			want: []FileStatus{
				{Path: "a.go", Staged: 'M', Unstaged: '.'},
				{Path: "b.go", Untracked: true, Staged: '?', Unstaged: '?'},
			},
		},
		{
			name: "empty",
			in:   "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePorcelainV2(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d entries, want %d: %+v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("entry %d = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
