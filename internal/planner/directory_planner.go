package planner

import (
	"fmt"
	"strings"

	"git-split/helpers"
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
	sortedDirs := helpers.SortMap(git.GroupFilesByDepthMap(files, p.Depth))
	var branches []plan.BranchPlan
	currentBase := p.Base
	for i, filesPerDir := range sortedDirs {
		branch := fmt.Sprintf("%s-split-%d", p.Target, i+1)
		op := plan.Operation{
			Type:    plan.OpApplyPath,
			Paths:   filesPerDir.Value,
			FromRef: p.Target,
		}
		branches = append(branches, plan.BranchPlan{
			Branch:        branch,
			Base:          currentBase,
			Operations:    []plan.Operation{op},
			Push:          p.Push,
			CreateMR:      p.CreateMR,
			MRTitle:       fmt.Sprintf("%s: Split %d", p.Target, i+1),
			MRDescription: fmt.Sprintf("This MR splits the changes in `%s` into a separate branch. Summary of files changed in this MR: \n - %s", filesPerDir.Key, strings.Join(filesPerDir.Value, "\n - ")),
		})
		currentBase = branch
	}
	return plan.Plan{Branches: branches}, nil
}
