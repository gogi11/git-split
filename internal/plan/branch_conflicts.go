package plan

import (
	"fmt"
	"git-split/helpers"
	"git-split/internal/git"
)

func FixBranchConflicts(p Plan, isPushing bool, autoDelete bool) error {
	var localConflicts []string
	var remoteConflicts []string
	for _, b := range p.Branches {
		if git.LocalBranchExists(b.Branch) {
			localConflicts = append(localConflicts, b.Branch)
		}
		if isPushing && git.RemoteBranchExists(p.Remote, b.Branch) {
			remoteConflicts = append(remoteConflicts, b.Branch)
		}
	}
	if len(localConflicts) == 0 && len(remoteConflicts) == 0 {
		return nil
	}

	fmt.Println("The following local and remote branches already exist:")
	if len(localConflicts) > 0 {
		fmt.Println("LOCAL:")
		for _, b := range localConflicts {
			fmt.Println(" -", b)
		}
	}
	if len(remoteConflicts) > 0 {
		fmt.Println("REMOTE:")
		for _, b := range remoteConflicts {
			fmt.Println(" -", b)
		}
	}

	if !autoDelete && !helpers.AskConfirmation("Delete them to continue?") {
		return fmt.Errorf("Aborting due to user selection.")
	}
	for _, b := range localConflicts {
		fmt.Println("Deleting branch:", b)
		err := git.DeleteLocalBranch(p.Remote, b)
		if err != nil {
			return fmt.Errorf("Failed to delete local branch: %s", b)
		}
	}
	for _, b := range remoteConflicts {
		fmt.Println("Deleting branch:", b)
		err := git.DeleteRemoteBranch(p.Remote, b)
		if err != nil {
			return fmt.Errorf("Failed to delete remote branch: %s", b)
		}
	}
	return nil
}
