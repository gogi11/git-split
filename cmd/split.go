package cmd

import (
	"fmt"
	"log"

	"git-split/internal/executor"
	"git-split/internal/plan"
	"git-split/internal/planner"

	"github.com/spf13/cobra"
)

var (
	base      string
	target    string
	size      int
	prefix    string
	push      bool
	createMR  bool
	dryRun    bool
	mode      string
	pathDepth int
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split commits into stacked branches",
	Run: func(cmd *cobra.Command, args []string) {
		var plannerImpl planner.Planner
		switch mode {
		case "directory":
			plannerImpl = planner.DirectoryPlanner{
				Base:     base,
				Target:   target,
				Depth:    pathDepth,
				Push:     push,
				CreateMR: createMR,
			}
		default:
			plannerImpl = planner.CommitPlanner{
				Base:     base,
				Target:   target,
				Size:     size,
				Push:     push,
				CreateMR: createMR,
			}
		}
		planning, err := plannerImpl.Build()
		if err != nil {
			log.Fatal(err)
		}
		plan.PrintPreview(planning)
		if dryRun {
			fmt.Println("Dry-run mode enabled. No changes were made.")
			return
		}
		executor.Execute(planning)
	},
}

func init() {

	splitCmd.Flags().StringVar(&base, "base", "main", "Base branch")
	splitCmd.Flags().StringVar(&target, "target", "", "Target branch")
	splitCmd.Flags().IntVar(&size, "size", 5, "Commits per branch")
	splitCmd.Flags().StringVar(&prefix, "prefix", "split", "Branch prefix")
	splitCmd.Flags().BoolVar(&push, "push", true, "Push branches")
	splitCmd.Flags().BoolVar(&createMR, "create-mr", false, "Create merge requests")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate actions (without actually pushing)")
	splitCmd.Flags().StringVar(&mode, "mode", "commit", "Spliting mode: commit | directory")
	splitCmd.Flags().IntVar(&pathDepth, "path-depth", 1, "Path depth for directory-based splitting")

	rootCmd.AddCommand(splitCmd)
}
