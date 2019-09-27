package cmd

import (
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/steps"
	"github.com/bdewey/git-stack/src/util"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the last run git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		runState := steps.LoadPreviousRunState()
		if runState == nil || runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to undo")
		}
		undoRunState := runState.CreateUndoRunState()
		steps.Run(&undoRunState)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(undoCmd)
}
