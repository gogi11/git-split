package plan

import (
	"fmt"
	"git-split/internal/mr"
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
	total := len(chunks)
	for i, chunk := range chunks {
		index := i + 1
		branch := fmt.Sprintf("%s-%d", prefix, index)
		description := mr.GenerateDescription(
			index,
			total,
			currentBase,
			branch,
			chunk,
		)
		branches = append(branches, BranchPlan{
			Branch:        branch,
			Base:          currentBase,
			Commits:       chunk,
			Push:          push,
			CreateMR:      createMR,
			MRTitle:       fmt.Sprintf("[%d/%d] %s", index, total, target),
			MRDescription: description,
		})
	}
	return Plan{
		Remote:   remote,
		Provider: provider,
		Repo:     repo,
		Branches: branches,
	}
}
