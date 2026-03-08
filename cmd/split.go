package cmd

import (
	"fmt"
	"log"

	"git-split/internal/git"
	"git-split/internal/mr"
	"git-split/internal/provider"

	"github.com/spf13/cobra"
)

var (
	base     string
	target   string
	size     int
	prefix   string
	push     bool
	createMR bool
	dryRun   bool
	remote   string
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split commits into stacked branches",
	Run: func(cmd *cobra.Command, args []string) {

		commits, err := git.GetCommitsBetween(base, target)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found %d commits\n", len(commits))

		chunks := git.ChunkCommits(commits, size)

		currentBase := base

		repoInfo, err := provider.ParseRemote(remote)
		if err != nil {
			log.Fatal(err)
		}

		for i, chunk := range chunks {

			branch := fmt.Sprintf("%s-%d", prefix, i+1)

			if dryRun {
				fmt.Printf("[DRY] create branch %s from %s\n", branch, currentBase)
				fmt.Printf("[DRY] cherry-pick %v\n", chunk)
			} else {

				err := git.CreateBranch(currentBase, branch)
				if err != nil {
					log.Fatal(err)
				}

				err = git.CherryPickCommits(chunk)
				if err != nil {
					log.Fatal(err)
				}
			}

			if push {

				if dryRun {
					fmt.Printf("[DRY] push %s\n", branch)
				} else {
					err := git.Push(remote, branch)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			if createMR {

				title := fmt.Sprintf("%s part %d", target, i+1)

				if dryRun {
					fmt.Printf("[DRY] create MR %s -> %s\n", currentBase, branch)
				} else {

					err := mr.Create(
						repoInfo,
						title,
						currentBase,
						branch,
					)

					if err != nil {
						log.Fatal(err)
					}
				}
			}

			currentBase = branch
		}
	},
}

func init() {

	splitCmd.Flags().StringVar(&base, "base", "", "Base branch")
	splitCmd.Flags().StringVar(&target, "target", "", "Target branch")
	splitCmd.Flags().IntVar(&size, "size", 5, "Commits per branch")
	splitCmd.Flags().StringVar(&prefix, "prefix", "split", "Branch prefix")
	splitCmd.Flags().BoolVar(&push, "push", false, "Push branches")
	splitCmd.Flags().BoolVar(&createMR, "create-mr", false, "Create merge requests")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions")
	splitCmd.Flags().StringVar(&remote, "remote", "origin", "Git remote")

	rootCmd.AddCommand(splitCmd)
}
