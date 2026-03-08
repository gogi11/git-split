package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-split",
	Short: "Split a large branch into smaller reviewable branches",
	Long:  "A CLI tool that splits a branch into multiple intermediate branches for easier merge request review.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
