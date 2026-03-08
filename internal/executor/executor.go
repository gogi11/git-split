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
		return err
	}
	repoInfo, err := provider.ParseRemote(remote)
	if err != nil {
		return err
	}

	for _, branch := range p.Branches {
		err = git.Checkout(branch.Base)
		if err != nil {
			return err
		}
		err := git.CreateBranch(branch.Base, branch.Branch)
		if err != nil {
			return err
		}
		for _, op := range branch.Operations {
			switch op.Type {
			case plan.OpCherryPick:
				err := git.CherryPickCommits(op.Commits)
				if err != nil {
					return err
				}
			case plan.OpApplyPath:
				for _, path := range op.Paths {
					err := git.ApplyPathFromBranch(
						op.FromRef,
						path,
					)
					if err != nil {
						return err
					}
				}
				msg := fmt.Sprintf("Apply paths from %s", op.FromRef)
				err := git.Commit(msg)
				if err != nil {
					return err
				}
			}
		}
		if branch.Push {
			err := git.Push(remote, branch.Branch)
			if err != nil {
				log.Fatal(err)
			}
			if branch.CreateMR {
				err := mr.Create(
					repoInfo,
					branch.MRTitle,
					branch.MRDescription,
					branch.Base,
					branch.Branch,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
