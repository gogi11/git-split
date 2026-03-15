package plan

import (
	"fmt"
	"git-split/helpers"
	"git-split/internal/git"
	"log"
)

func FixBranchConflicts(p Plan, autoDelete bool) error {
	var existing []string
	for _, b := range p.Branches {
		if git.LocalBranchExists(b.Branch) || git.RemoteBranchExists(p.Remote, b.Branch) {
			existing = append(existing, b.Branch)
		}
	}

	if len(existing) > 0 {
		fmt.Println("The following branches already exist:")
		for _, b := range existing {
			fmt.Println(" -", b)
		}
		if !autoDelete && !helpers.AskConfirmation("Delete them to continue?") {
			log.Fatal("Aborting.")
			return fmt.Errorf("branches exist already")
		}
		for _, b := range existing {
			fmt.Println("Deleting branch:", b)
			err := git.DeleteBranch(p.Remote, b)
			if err != nil {
				log.Fatal("Failed to delete branch:", b)
				return err
			}
		}
	}
	return nil
}
