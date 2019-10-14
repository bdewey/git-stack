package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/bdewey/git-stack/src/script"

	"github.com/bdewey/git-stack/src/git"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Move to the next branch in the stack",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		moveToChildBranch()
	},
}

func moveToChildBranch() {
	currentBranch := git.GetCurrentBranchName()
	childBranches := git.GetChildBranches(currentBranch)
	switch len(childBranches) {
	case 0:
		fmt.Println("At the top of the stack, nowhere to go!")
		os.Exit(1)
	case 1:
		script.RunCommand("git", "checkout", childBranches[0])
	default:
		fmt.Println("More than one child branch; don't know where to go!")
		fmt.Println("Options:")
		sort.Strings(childBranches)
		for _, childBranch := range childBranches {
			fmt.Println(" -", childBranch)
		}
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(upCmd)
}
