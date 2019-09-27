package cmd

import (
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/prompt"
	"github.com/bdewey/git-stack/src/script"
	"github.com/bdewey/git-stack/src/steps"
	"github.com/bdewey/git-stack/src/util"

	"github.com/spf13/cobra"
)

type appendConfig struct {
	ParentBranch string
	TargetBranch string
}

var appendCommand = &cobra.Command{
	Use:   "append <branch>",
	Short: "Creates a new feature branch as a child of the current branch",
	Long: `Creates a new feature branch as a direct child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository
if and only if new-branch-push-flag is true,
and brings over all uncommitted changes to the new feature branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getAppendConfig(args)
		stepList := getAppendStepList(config)
		runState := steps.NewRunState("append", stepList)
		steps.Run(runState)
	},
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func getAppendConfig(args []string) (result appendConfig) {
	result.ParentBranch = git.GetCurrentBranchName()
	result.TargetBranch = args[0]
	if git.HasRemote("origin") && !git.IsOffline() {
		script.Fetch()
	}
	git.EnsureDoesNotHaveBranch(result.TargetBranch)
	prompt.EnsureKnowsParentBranches([]string{result.ParentBranch})
	return
}

func init() {
	RootCmd.AddCommand(appendCommand)
}
