package executor

import (
	"fmt"

	"git-split/internal/filechanges"
	"git-split/internal/git"
	"git-split/internal/mr"
	"git-split/internal/plan"
)

func Execute(p plan.Plan) error {
	fmt.Printf("Executing plan for repository: %s\n", p.Repo.Name)
	for _, branch := range p.Branches {
		err := git.Checkout(branch.Base)
		if err != nil {
			return err
		}
		fmt.Printf("Checked out to branch: %s\n", branch.Base)
		err = git.CreateBranch(branch.Base, branch.Branch)
		if err != nil {
			return err
		}
		fmt.Printf("Created branch: %s\n", branch.Branch)
		for _, op := range branch.Operations {
			switch op.Type {
			case plan.OpCherryPick:
				err := git.CherryPickCommits(op.Commits)
				if err != nil {
					return err
				}
			case plan.OpApplyPath:
				for _, fc := range op.FileChanges {
					switch fc.Action {
					case filechanges.MODIFIED, filechanges.ADDED:
						git.ApplyPathFromBranch(op.FromRef, fc.Path) // checkout/update file
					case filechanges.DELETED:
						git.DeleteFile(fc.Path) // remove file
					case filechanges.RENAMED:
						git.MoveFile(fc.OldPath, fc.Path) // optional: rename
					}
				}
				err := git.Commit(branch.MRTitle)
				if err != nil {
					return err
				}
			}
		}
		if branch.Push {
			err := git.Push(p.Remote, branch.Branch)
			if err != nil {
				return err
			}
			fmt.Printf("Pushed branch: %s\n", branch.Branch)
			if branch.CreateMR {
				err := mr.Create(
					p.Repo,
					branch.MRTitle,
					branch.MRDescription,
					branch.Base,
					branch.Branch,
				)
				if err != nil {
					return err
				}
			}
			fmt.Printf("Created MR for branch: %s\n", branch.Branch)
		}
	}
	return nil
}
