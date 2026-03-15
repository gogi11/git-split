package cmd

import (
	"fmt"
	"log"

	"git-split/internal/git"
	"git-split/internal/planner"
)

func LoadRepo(target *string, base *string) error {
	git.Fetch()
	if *target == "" {
		current, err := git.GetCurrentBranch()
		if err != nil {
			log.Fatal("Unable to detect current branch")
			return fmt.Errorf("unable to detect current branch: %w", err)
		}

		*target = current
		fmt.Printf("Target branch not specified, using current branch: %s\n", *target)
	}

	if *base == *target {
		log.Fatal("Base and target branches cannot be the same")
		return fmt.Errorf("base and target branches cannot be the same")
	}
	err := git.RebaseOnto(*base)
	if err != nil {
		log.Fatal("Rebase failed, aborting split.")
		return err
	}

	err = git.ForcePush(*target, *base)
	if err != nil {
		log.Fatal("Failed to push rebased branch.")
		return err
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
