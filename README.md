# zinc

A terminal UI for git, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
and [Lipgloss](https://github.com/charmbracelet/lipgloss). Think
[lazygit](https://github.com/jesseduffield/lazygit): a lazygit-style
sidebar/main-panel layout over the git CLI, driven entirely from the
keyboard.

This is the core git experience — status, staging, diffs, commits,
branches, log, and stash. GitHub CLI integration for pull requests and
Actions workflows is planned as a follow-up.

## Install

```sh
go install github.com/jamesjohnsdev/zinc/cmd/zinc@latest
```

Requires the `git` binary on `PATH`. Run `zinc` from inside (or below) any
git repository.

## Layout

```
┌─ Files ────────┐┌─ Diff ─────────────────┐
│                ││                        │
├─ Branches ─────┤│                        │
│                ││                        │
├─ Commits ──────┤│                        │
│                ││                        │
├─ Stash ────────┤│                        │
└────────────────┘└────────────────────────┘
```

`tab` / `shift+tab` cycle focus between the four sidebar panels. The Diff
panel on the right always reflects the file currently selected in Files.

## Keybindings

**Global**

| Key       | Action              |
| --------- | ------------------- |
| `tab`     | switch panel focus   |
| `shift+tab` | switch panel focus (reverse) |
| `↑`/`k`, `↓`/`j` | move selection in the focused panel |
| `r`       | refresh all panels   |
| `q`, `ctrl+c` | quit             |

**Files**

| Key   | Action                              |
| ----- | ------------------------------------ |
| `space` | stage/unstage the selected file    |
| `a`   | stage all changes                    |
| `c`   | open the commit message prompt       |
| `s`   | stash all changes (including untracked) |
| `D`   | discard the selected file's changes (asks to confirm) |
| `ctrl+u` / `ctrl+d` | scroll the diff panel   |

**Branches**

| Key     | Action                                          |
| ------- | ------------------------------------------------ |
| `enter` | check out the selected branch (tracks it locally first if remote) |
| `n`     | prompt for a new branch name, create and check it out |
| `d`     | delete the selected local branch (asks to confirm) |

**Stash**

| Key | Action                          |
| --- | -------------------------------- |
| `p` | pop the selected stash entry     |
| `d` | drop the selected stash entry (asks to confirm) |

Prompts (commit message, new branch name) accept `enter` to submit and
`esc` to cancel. Confirmation dialogs accept `y` to proceed and `n`/`esc`
to cancel.

## Development

```sh
go build ./...
go test ./...
golangci-lint run
```

CI runs golangci-lint, govulncheck, `go build`, and `go test -race` on
every push and pull request; see `.github/workflows/`.

## License

[MIT](LICENSE)
