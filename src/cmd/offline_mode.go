package cmd

import (
	"github.com/bdewey/git-stack/src/cfmt"
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/util"
	"github.com/spf13/cobra"
)

var offlineCommand = &cobra.Command{
	Use:   "offline [(true | false)]",
	Short: "Displays or sets offline mode",
	Long: `Displays or sets offline mode

Use offline mode to maximize Git Town capabilities when working in an offline environment, such as on an airplane.
When Git Town is in offline mode, certain commands will lose networking and upstream syncing capabilities but will
continue to work locally. Commands that rely on network functionality will be completely disabled.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printOfflineFlag()
		} else {
			setOfflineFlag(util.StringToBool(args[0]))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return validateBooleanArgument(args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
}

func printOfflineFlag() {
	cfmt.Println(git.GetPrintableOfflineFlag())
}

func setOfflineFlag(value bool) {
	git.UpdateOffline(value)
}

func init() {
	RootCmd.AddCommand(offlineCommand)
}
