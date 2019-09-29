package cmd

import (
	"fmt"
	"os"

	"github.com/bdewey/git-stack/src/command"
	"github.com/bdewey/git-stack/src/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra object
var RootCmd = &cobra.Command{
	Use:   "git-stack",
	Short: "Stacked pull request support for Github",
	Long: `Git Stack simplifies the process of stacking pull requests for review on Github.

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		command.SetDebug(debugFlag)
	},
}

// Execute runs the Cobra stack
func Execute() {
	git.EnsureVersionRequirementSatisfied()
	color.NoColor = false // Prevent color from auto disable

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Developer tool to print git commands run under the hood")
}
