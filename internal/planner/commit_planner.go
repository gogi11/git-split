package planner

import (
	"fmt"
	"strings"

	"git-split/internal/git"
	"git-split/internal/plan"
)

type CommitPlanner struct {
	Base     string
	Target   string
	Size     int
	Push     bool
	CreateMR bool
}

func (p CommitPlanner) Build(remote string) (plan.Plan, error) {
	resultingPlan, err := InitializePlan(remote)
	if err != nil {
		return plan.Plan{}, err
	}
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
			Branch:        branch,
			Base:          currentBase,
			Operations:    []plan.Operation{op},
			Push:          p.Push,
			CreateMR:      p.CreateMR,
			MRTitle:       fmt.Sprintf("%s: Split %d", p.Target, i+1),
			MRDescription: fmt.Sprintf("This MR splits the commits in `%s` into a separate branch. Hashes of commits in this MR: \n - %s", p.Target, strings.Join(chunk, "\n - ")),
		})
		currentBase = branch
	}
	resultingPlan.Branches = branches
	return resultingPlan, nil
}
