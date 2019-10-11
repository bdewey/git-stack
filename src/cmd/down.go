package cmd

import (
	"fmt"
	"os"

	"github.com/bdewey/git-stack/src/script"

	"github.com/bdewey/git-stack/src/git"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Move to the previous branch in the stack",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		moveToParentBranch()
	},
}

func moveToParentBranch() {
	currentBranch := git.GetCurrentBranchName()
	parentBranch := git.GetParentBranch(currentBranch)
	if parentBranch != "" {
		script.RunCommand("git", "checkout", parentBranch)
	} else {
		fmt.Println("No parent branch in the stack.")
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(downCmd)
}
