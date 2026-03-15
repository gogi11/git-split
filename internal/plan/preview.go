package plan

import (
	"fmt"
)

func PrintPreview(p Plan) {
	fmt.Println("")
	fmt.Println("========== STACK PLAN ==========")
	fmt.Printf("Remote:   %s\n", p.Remote)
	fmt.Printf("Provider: %s\n", p.Repo.Provider)
	fmt.Printf("Repo:     %s\n", p.Repo.Name)
	fmt.Println("")
	for i, b := range p.Branches {
		fmt.Printf("[%d] Branch: %s\n", i+1, b.Branch)
		fmt.Printf("    Base: %s\n", b.Base)
		for _, op := range b.Operations {
			switch op.Type {
			case OpCherryPick:
				fmt.Printf("    Cherry-pick Commits (%d):\n", len(op.Commits))
				for _, c := range op.Commits {
					fmt.Printf("      - %s\n", c)
				}
			case OpApplyPath:
				fmt.Printf("    Apply Paths from %s:\n", op.FromRef)
				for _, fileChange := range op.FileChanges {
					fmt.Printf("      - %s %s\n", fileChange.Action, fileChange.Path)
				}
			}
		}
		if b.Push {
			fmt.Println("    Action: push branch")
		}
		if b.CreateMR {
			fmt.Printf("    MR: %s -> %s\n", b.Base, b.Branch)
			fmt.Printf("    Title: %s\n", b.MRTitle)
		}
		fmt.Println("")
	}
	fmt.Println("================================")
	fmt.Println("")
}
