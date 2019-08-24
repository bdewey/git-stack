package cmd

import (
	"fmt"
	"unicode/utf8"

	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"strings"
)

type branchInfo struct {
	Branch    string
	IsCurrent bool
	PrInfo    string
}

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Display the current working stack",
	Long: `Display the current working stack

Shows information about the current stack of changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		printBranchInfo(getBranchInfo())
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

func getBranchInfo() []branchInfo {
	var info []branchInfo
	currentBranch := git.GetCurrentBranchName()
	for _, branch := range git.GetAncestorBranches(currentBranch) {
		info = append(info, branchInfo{Branch: branch})
	}
	currentChildren := []string{currentBranch}
	for ; len(currentChildren) == 1; currentChildren = git.GetChildBranches(currentChildren[0]) {
		isCurrent := currentChildren[0] == currentBranch
		info = append(info, branchInfo{Branch: currentChildren[0], IsCurrent: isCurrent})
	}

	// Look for PR info for each of the branches.
	driver := drivers.GetActiveDriver()
	for i := range info {
		if i > 0 {
			canMerge, defaultCommitMessage, _ := driver.CanMergePullRequest(info[i].Branch, info[i-1].Branch)
			if canMerge {
				info[i].PrInfo = defaultCommitMessage
			}
		}
	}
	return info
}

func printBranchInfo(branchInfo []branchInfo) {
	longestBranchRuneCount := utf8.RuneCountInString("Branch")
	longestPrString := 0
	for _, info := range branchInfo {
		currentRuneCount := utf8.RuneCountInString(info.Branch)
		// The current branch needs to display "* " at the start
		// TODO: Factor "* " into a constant where I can find its rune length
		if info.IsCurrent {
			currentRuneCount += 2
		}
		// Always allow for 2 spaces after the end of the string
		currentRuneCount += 2
		if currentRuneCount > longestBranchRuneCount {
			longestBranchRuneCount = currentRuneCount
		}
		currentPrStringCount := utf8.RuneCountInString(info.PrInfo)
		if currentPrStringCount > longestPrString {
			longestPrString = currentPrStringCount
		}
	}

	fmt.Print("Branch")
	fmt.Print(strings.Repeat(" ", longestBranchRuneCount-utf8.RuneCountInString("Branch")))
	if longestPrString > 0 {
		fmt.Print("Github PR")
	}
	fmt.Println("")
	fmt.Print(strings.Repeat("=", longestBranchRuneCount-2))
	fmt.Print("  ")
	fmt.Println(strings.Repeat("=", longestPrString))

	for _, info := range branchInfo {
		if info.IsCurrent {
			cfmt.Print("* ")
		}
		branchColor := color.New(color.Bold)
		if info.IsCurrent {
			branchColor = branchColor.Add(color.FgCyan)
		}
		cfmt.Print(branchColor.Sprintf(info.Branch))
		spacesToPrint := longestBranchRuneCount - utf8.RuneCountInString(info.Branch)
		if info.IsCurrent {
			spacesToPrint -= 2
		}
		cfmt.Print(strings.Repeat(" ", spacesToPrint))
		cfmt.Println(info.PrInfo)
	}
}

func init() {
	RootCmd.AddCommand(stackCmd)
}
