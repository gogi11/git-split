package executor

import (
	"fmt"
	"log"

	"git-split/internal/git"
	"git-split/internal/mr"
	"git-split/internal/plan"
	"git-split/internal/provider"
)

func Execute(p plan.Plan) error {
	remote, err := git.GetRemoteURL()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Remote URL: %s\n", remote)
	repoInfo, err := provider.ParseRemote(remote)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Executing plan for repository: %s\n", repoInfo)

	for _, branch := range p.Branches {
		err = git.Checkout(branch.Base)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("Checked out to branch: %s\n", branch.Base)
		err := git.CreateBranch(branch.Base, branch.Branch)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Printf("Created branch: %s\n", branch.Branch)
		for _, op := range branch.Operations {
			switch op.Type {
			case plan.OpCherryPick:
				err := git.CherryPickCommits(op.Commits)
				if err != nil {
					log.Fatal(err)
					return err
				}
			case plan.OpApplyPath:
				for _, path := range op.Paths {
					err := git.ApplyPathFromBranch(
						op.FromRef,
						path,
					)
					if err != nil {
						log.Fatal(err)
						return err
					}
				}
				fmt.Printf("Applied these files: %s\n", op.Paths)
				msg := fmt.Sprintf("Apply paths from %s", op.FromRef)
				err := git.Commit(msg)
				if err != nil {
					log.Fatal(err)
					return err
				}
			}
		}
		if branch.Push {
			err := git.Push(remote, branch.Branch)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Pushed branch: %s\n", branch.Branch)
			if branch.CreateMR {
				err := mr.Create(
					repoInfo,
					branch.MRTitle,
					branch.MRDescription,
					branch.Base,
					branch.Branch,
				)
				if err != nil {
					log.Fatal(err)
					return err
				}
			}
			fmt.Printf("Created MR for branch: %s\n", branch.Branch)
		}
	}
	return nil
}
