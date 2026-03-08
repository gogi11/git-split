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
	Remote   string
	Provider string
	Repo     string
	Push     bool
	CreateMR bool
}

func (p DirectoryPlanner) Build() (plan.Plan, error) {

	files, err := git.GetChangedFiles(p.Base, p.Target)
	if err != nil {
		return plan.Plan{}, err
	}

	groups := git.GroupFilesByDirectory(files)

	var chunks [][]string

	for dir, files := range groups {

		fmt.Println("Directory:", dir)

		commits, err := git.GetCommitsForFiles(
			p.Base,
			p.Target,
			files,
		)

		if err != nil {
			return plan.Plan{}, err
		}

		chunks = append(chunks, commits)
	}

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
