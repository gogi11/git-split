package planner

import (
	"fmt"

	"git-split/internal/git"
	"git-split/internal/plan"
)

type DirectoryPlanner struct {
	Base     string
	Target   string
	Prefix   string
	Depth    int
	Push     bool
	CreateMR bool
}

func (p DirectoryPlanner) Build() (plan.Plan, error) {

	files, err := git.GetChangedFiles(p.Base, p.Target)
	if err != nil {
		return plan.Plan{}, err
	}
	dirs := git.GroupFilesByDepth(files, p.Depth)
	var branches []plan.BranchPlan
	currentBase := p.Base
	var accumulated []string
	for i, dir := range dirs {
		accumulated = append(accumulated, dir)
		branch := fmt.Sprintf("%s-split-%d", p.Target, i+1)
		op := plan.Operation{
			Type:    plan.OpApplyPath,
			Paths:   []string{dir},
			FromRef: p.Target,
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
