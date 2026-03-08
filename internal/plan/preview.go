package plan

import (
	"fmt"
)

func PrintPreview(p Plan) {

	fmt.Println("")
	fmt.Println("========== STACK PLAN ==========")

	fmt.Printf("Remote:   %s\n", p.Remote)
	fmt.Printf("Provider: %s\n", p.Provider)
	fmt.Printf("Repo:     %s\n", p.Repo)

	fmt.Println("")

	for i, b := range p.Branches {

		fmt.Printf("[%d] Branch: %s\n", i+1, b.Branch)
		fmt.Printf("    Base: %s\n", b.Base)

		fmt.Printf("    Commits (%d):\n", len(b.Commits))

		for _, c := range b.Commits {
			fmt.Printf("      - %s\n", c)
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
