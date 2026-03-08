package cmd

import (
	"fmt"
	"log"

	"git-split/internal/git"

	"github.com/spf13/cobra"
)

var (
	base   string
	target string
	size   int
	prefix string
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split commits into multiple branches",
	Run: func(cmd *cobra.Command, args []string) {

		if base == "" || target == "" {
			log.Fatal("--base and --target are required")
		}

		commits, err := git.GetCommitsBetween(base, target)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found %d commits\n", len(commits))

		chunks := git.ChunkCommits(commits, size)

		currentBase := base

		for i, chunk := range chunks {

			branchName := fmt.Sprintf("%s-%d", prefix, i+1)

			fmt.Printf("\nCreating branch %s\n", branchName)

			err := git.CreateBranch(currentBase, branchName)
			if err != nil {
				log.Fatal(err)
			}

			err = git.CherryPickCommits(chunk)
			if err != nil {
				log.Fatal(err)
			}

			currentBase = branchName
		}

		fmt.Println("\nFinished splitting branches.")
	},
}

func init() {

	splitCmd.Flags().StringVar(&base, "base", "", "Base branch")
	splitCmd.Flags().StringVar(&target, "target", "", "Target branch")
	splitCmd.Flags().IntVar(&size, "size", 5, "Number of commits per branch")
	splitCmd.Flags().StringVar(&prefix, "prefix", "split", "Branch name prefix")

	rootCmd.AddCommand(splitCmd)
}
