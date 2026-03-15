package cmd

import (
	"fmt"

	"git-split/internal/git"
	"git-split/internal/planner"
)

func LoadRepo(target *string, base *string, autoDelete bool) error {
	if !git.WorkingTreeClean() {
		return fmt.Errorf("Working tree is not clean. Please commit or stash your changes before running the split.")
	}
	git.Fetch(autoDelete)
	if *target == "" {
		current, err := git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("Unable to detect current branch: %w", err)
		}

		*target = current
		fmt.Printf("Target branch not specified, using current branch: %s\n", *target)
	}

	if *base == *target {
		return fmt.Errorf("Base and target branches cannot be the same")
	}
	err := git.RebaseOnto(*base)
	if err != nil {
		return fmt.Errorf("Rebase failed, aborting split: %w", err)
	}

	err = git.ForcePush(*target, *base)
	if err != nil {
		return fmt.Errorf("Failed to push rebased branch: %w", err)
	}
	return nil
}

func SelectPlanner(mode string) planner.Planner {
	var plannerImpl planner.Planner
	switch mode {
	case "directory":
		plannerImpl = planner.DirectoryPlanner{
			Base:     base,
			Target:   target,
			Depth:    pathDepth,
			Push:     push,
			CreateMR: createMR,
		}
	default:
		plannerImpl = planner.CommitPlanner{
			Base:     base,
			Target:   target,
			Size:     size,
			Push:     push,
			CreateMR: createMR,
		}
	}
	return plannerImpl
}
