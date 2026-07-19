package gh

import (
	"context"
	"encoding/json"
	"fmt"
)

// Repo is summary information about a GitHub repository.
type Repo struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	Owner         Author `json:"owner"`
	Description   string `json:"description"`
	DefaultBranch struct {
		Name string `json:"name"`
	} `json:"defaultBranchRef"`
	IsPrivate      bool   `json:"isPrivate"`
	StargazerCount int    `json:"stargazerCount"`
	URL            string `json:"url"`
}

const repoViewFields = "name,nameWithOwner,owner,description,defaultBranchRef,isPrivate,stargazerCount,url"

// ViewRepo returns summary information for repo, or the repository
// resolved from the working directory's git remote if repo is empty.
//
// Unlike pr/issue/run subcommands, `gh repo view` takes the repository as a
// positional argument rather than a --repo flag.
func (r *Runner) ViewRepo(ctx context.Context, repo string) (Repo, error) {
	out, err := r.run(ctx, viewRepoArgs(repo)...)
	if err != nil {
		return Repo{}, err
	}

	var rp Repo
	if err := json.Unmarshal([]byte(out), &rp); err != nil {
		return Repo{}, fmt.Errorf("parse repo view: %w", err)
	}

	return rp, nil
}

// viewRepoArgs builds the argument list for `gh repo view`. Unlike pr/
// issue/run subcommands, it takes the repository as a positional argument
// rather than a --repo flag.
func viewRepoArgs(repo string) []string {
	args := []string{"repo", "view"}
	if repo != "" {
		args = append(args, repo)
	}
	return append(args, "--json", repoViewFields)
}
