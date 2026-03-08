package planner

import (
	"fmt"

	"git-split/internal/git"
	"git-split/internal/plan"
)

type CommitPlanner struct {
	Base     string
	Target   string
	Size     int
	Prefix   string
	Push     bool
	CreateMR bool
}

func (p CommitPlanner) Build() (plan.Plan, error) {
	commits, err := git.GetCommitsBetween(p.Base, p.Target)
	if err != nil {
		return plan.Plan{}, err
	}
	chunks := git.ChunkCommits(commits, p.Size)
	var branches []plan.BranchPlan
	currentBase := p.Base
	for i, chunk := range chunks {
		branch := fmt.Sprintf("%s-split-%d", p.Target, i+1)
		op := plan.Operation{
			Type:    plan.OpCherryPick,
			Commits: chunk,
		}
		branches = append(branches, plan.BranchPlan{
			Branch:     branch,
			Base:       currentBase,
			Operations: []plan.Operation{op},
			Push:       p.Push,
			CreateMR:   p.CreateMR,
		})
		currentBase = branch
	}
	return plan.Plan{Branches: branches}, nil
}
