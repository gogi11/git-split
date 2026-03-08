package planner

import (
	"git-split/internal/git"
	"git-split/internal/plan"
)

type CommitPlanner struct {
	Base     string
	Target   string
	Size     int
	Prefix   string
	Remote   string
	Provider string
	Repo     string
	Push     bool
	CreateMR bool
}

func (p CommitPlanner) Build() (plan.Plan, error) {

	commits, err := git.GetCommitsBetween(p.Base, p.Target)
	if err != nil {
		return plan.Plan{}, err
	}

	chunks := git.ChunkCommits(commits, p.Size)

	return plan.BuildPlan(
		p.Remote,
		p.Provider,
		p.Repo,
		p.Base,
		p.Target,
		p.Prefix,
		chunks,
		p.Push,
		p.CreateMR,
	), nil
}
