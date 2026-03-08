package cmd

import (
	"fmt"
	"log"

	"git-split/internal/git"
	"git-split/internal/mr"
	"git-split/internal/plan"
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

		remote, err := git.GetRemoteURL()
		if err != nil {
			log.Fatal(err)
		}
		repoInfo, err := provider.ParseRemote(remote)
		if err != nil {
			log.Fatal(err)
		}

		planning := plan.BuildPlan(
			remote,
			string(repoInfo.Provider),
			repoInfo.Owner+"/"+repoInfo.Name,
			base,
			target,
			prefix,
			chunks,
			push,
			createMR,
		)
		plan.PrintPreview(planning)
		if dryRun {
			fmt.Println("Dry-run mode enabled. No changes were made.")
			return
		}
		for i, branchPlan := range planning.Branches {

			branch := fmt.Sprintf("%s-%d", prefix, i+1)
			err := git.CreateBranch(currentBase, branch)
			if err != nil {
				log.Fatal(err)
			}
			err = git.CherryPickCommits(branchPlan.Commits)
			if err != nil {
				log.Fatal(err)
			}
			if push {
				err := git.Push(remote, branch)
				if err != nil {
					log.Fatal(err)
				}
			}
			if createMR {
				err := mr.Create(
					repoInfo,
					branchPlan.MRTitle,
					branchPlan.MRDescription,
					currentBase,
					branch,
				)

				if err != nil {
					log.Fatal(err)
				}
			}
			currentBase = branch
		}
	},
}

func init() {

	splitCmd.Flags().StringVar(&base, "base", "main", "Base branch")
	splitCmd.Flags().StringVar(&target, "target", "", "Target branch")
	splitCmd.Flags().IntVar(&size, "size", 5, "Commits per branch")
	splitCmd.Flags().StringVar(&prefix, "prefix", "split", "Branch prefix")
	splitCmd.Flags().BoolVar(&push, "push", false, "Push branches")
	splitCmd.Flags().BoolVar(&createMR, "create-mr", false, "Create merge requests")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions")

	rootCmd.AddCommand(splitCmd)
}
