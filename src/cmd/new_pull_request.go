package cmd

import (
	"github.com/bdewey/git-stack/src/drivers"
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/prompt"
	"github.com/bdewey/git-stack/src/script"
	"github.com/bdewey/git-stack/src/steps"
	"github.com/bdewey/git-stack/src/util"
	"github.com/spf13/cobra"
)

type newPullRequestConfig struct {
	InitialBranch  string
	BranchesToSync []string
}

var newPullRequestCommand = &cobra.Command{
	Use:   "new-pull-request",
	Short: "Creates a new pull request",
	Long: `Creates a new pull request

Syncs the current branch
and opens a browser window to the new pull request page of your repository.

The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config git-town.code-hosting-driver <driver>"
where driver is "github", "gitlab", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config git-town.code-hosting-origin-hostname <hostname>"
where hostname matches what is in your ssh config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getNewPullRequestConfig()
		stepList := getNewPullRequestStepList(config)
		runState := steps.NewRunState("new-pull-request", stepList)
		steps.Run(runState)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
			git.ValidateIsOnline,
			drivers.ValidateHasDriver,
		)
	},
}

func getNewPullRequestConfig() (result newPullRequestConfig) {
	if git.HasRemote("origin") {
		script.Fetch()
	}
	result.InitialBranch = git.GetCurrentBranchName()
	prompt.EnsureKnowsParentBranches([]string{result.InitialBranch})
	result.BranchesToSync = append(git.GetAncestorBranches(result.InitialBranch), result.InitialBranch)
	return
}

func getNewPullRequestStepList(config newPullRequestConfig) (result steps.StepList) {
	// for _, branchName := range config.BranchesToSync {
	// 	result.AppendList(steps.GetSyncBranchSteps(branchName, true))
	// }
	result.Wrap(steps.WrapOptions{RunInGitRoot: true, StashOpenChanges: true})
	result.Append(&steps.CreatePullRequestStep{BranchName: config.InitialBranch})
	return
}

func init() {
	RootCmd.AddCommand(newPullRequestCommand)
}
