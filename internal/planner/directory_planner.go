package planner

import (
	"fmt"

	"git-split/helpers"
	"git-split/internal/filechanges"
	"git-split/internal/git"
	"git-split/internal/plan"
)

type DirectoryPlanner struct {
	Base     string
	Target   string
	Depth    int
	Push     bool
	CreateMR bool
}

func (p DirectoryPlanner) Build(remote string) (plan.Plan, error) {
	resultingPlan, err := InitializePlan(remote)
	if err != nil {
		return plan.Plan{}, err
	}
	actions, paths, err := git.GetChangedFilesWithStatus(p.Base, p.Target)
	if err != nil {
		return plan.Plan{}, err
	}
	files, err := filechanges.ConvertFileWithStatusLinesToFileChange(actions, paths)
	if err != nil {
		return plan.Plan{}, err
	}
	sortedDirs := helpers.SortMap(filechanges.GroupFilesByDepthMap(files, p.Depth))
	var branches []plan.BranchPlan
	currentBase := p.Base
	for i, filesPerDir := range sortedDirs {
		branch := fmt.Sprintf("%s-split-%d", p.Target, i+1)
		op := plan.Operation{
			Type:        plan.OpApplyPath,
			FileChanges: filesPerDir.Value,
			FromRef:     p.Target,
		}

		title := fmt.Sprintf("%s: Split %d", p.Target, i+1)

		description := fmt.Sprintf("This MR splits the changes in `%s` into a separate branch. Summary of files changed in this MR:", filesPerDir.Key)
		for _, f := range filesPerDir.Value {
			description += fmt.Sprintf("\n - %s %s", f.Action, f.Path)
		}

		branches = append(branches, plan.BranchPlan{
			Branch:        branch,
			Base:          currentBase,
			Operations:    []plan.Operation{op},
			Push:          p.Push,
			CreateMR:      p.CreateMR,
			MRTitle:       title,
			MRDescription: description,
		})
		currentBase = branch
	}
	resultingPlan.Branches = branches
	return resultingPlan, nil
}
