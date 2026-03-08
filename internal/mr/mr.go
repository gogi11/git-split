package mr

import (
	"fmt"

	"git-split/internal/provider"
)

func Create(repo provider.Repo, title, description, base, head string) error {
	switch repo.Provider {
	case provider.GitHub:
		return provider.CreateGitHubPR(repo, title, description, base, head)
	case provider.GitLab:
		return provider.CreateGitLabMR(repo, title, description, base, head)
	case provider.GitBucket:
		return provider.CreateGitBucketPR(repo, title, description, base, head)
	}
	return fmt.Errorf("unknown provider")
}
