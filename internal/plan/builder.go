package plan

import (
	"fmt"
)

func BuildPlan(
	remote string,
	provider string,
	repo string,
	base string,
	target string,
	prefix string,
	chunks [][]string,
	push bool,
	createMR bool,
) Plan {

	var branches []BranchPlan

	currentBase := base

	for i, chunk := range chunks {

		branch := fmt.Sprintf("%s-%d", prefix, i+1)

		branches = append(branches, BranchPlan{
			Branch:   branch,
			Base:     currentBase,
			Commits:  chunk,
			Push:     push,
			CreateMR: createMR,
			MRTitle:  fmt.Sprintf("%s part %d", target, i+1),
		})

		currentBase = branch
	}

	return Plan{
		Remote:   remote,
		Provider: provider,
		Repo:     repo,
		Branches: branches,
	}
}
