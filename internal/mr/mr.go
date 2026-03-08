package mr

import (
	"fmt"

	"git-split/internal/provider"
)

func Create(repo provider.Repo, title, base, head string) error {

	switch repo.Provider {
	case provider.GitHub:
		return provider.CreateGitHubPR(repo, title, base, head)
	case provider.GitLab:
		return provider.CreateGitLabMR(repo, title, base, head)
	case provider.GitBucket:
		return provider.CreateGitBucketPR(repo, title, base, head)
	}
	return fmt.Errorf("unknown provider")
}
