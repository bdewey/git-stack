package cmd

import (
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
	script.RunCommand("git", "checkout", parentBranch)
}

func init() {
	RootCmd.AddCommand(downCmd)
}
