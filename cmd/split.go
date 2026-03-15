package cmd

import (
	"fmt"
	"log"

	"git-split/internal/executor"
	"git-split/internal/git"
	"git-split/internal/plan"

	"github.com/spf13/cobra"
)

var (
	base       string
	target     string
	size       int
	push       bool
	createMR   bool
	dryRun     bool
	mode       string
	pathDepth  int
	autoDelete bool
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split commits into stacked branches",
	Run: func(cmd *cobra.Command, args []string) {
		err := LoadRepo(&target, &base, autoDelete)
		if err != nil {
			log.Fatal(err)
		}
		remote, err := git.SelectRemote()
		if err != nil {
			log.Fatal(err)
		}
		plannerImpl := SelectPlanner(mode)
		planning, err := plannerImpl.Build(remote)
		if err != nil {
			log.Fatal(err)
		}
		plan.PrintPreview(planning)
		if dryRun {
			fmt.Println("Dry-run mode enabled. No changes were made.")
			return
		}
		err = plan.FixBranchConflicts(planning, push, autoDelete)
		if err != nil {
			log.Fatal(err)
		}
		err = executor.Execute(planning)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	splitCmd.Flags().StringVar(&base, "base", "main", "Base branch")
	splitCmd.Flags().StringVar(&target, "target", "", "Target branch")
	splitCmd.Flags().IntVar(&size, "size", 5, "Commits per branch")
	splitCmd.Flags().BoolVar(&push, "push", true, "Push branches")
	splitCmd.Flags().BoolVar(&createMR, "create-mr", false, "Create merge requests")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions (without actually pushing)")
	splitCmd.Flags().StringVar(&mode, "mode", "commit", "Spliting mode: commit | directory")
	splitCmd.Flags().IntVar(&pathDepth, "depth", 2, "Path depth for directory-based splitting")
	splitCmd.Flags().BoolVar(&autoDelete, "delete", false, "Sets delete on everything automatically (local/remote, fetch prune, etc.)")

	rootCmd.AddCommand(splitCmd)
}
