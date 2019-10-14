package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bdewey/git-stack/src/drivers"
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/script"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// branchInfo contains all of the information needed to display a branch.
type branchInfo struct {
	// Branch is the name of the branch.
	Branch string
	// IsCurrent is true if this is the user's current branch.
	IsCurrent bool
	// If it exists, the PR number. 0 otherwise.
	PrNumber int
	// If there is a PR, the merge state. Empty otherwise.
	PrState string
	// git log --oneline output of the commits that are unique to this branch
	BranchCommits []string
}

var stackCmd = &cobra.Command{
	Use:   "list",
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

// getBranchInfo returns a slice with branchInfo structures for all branches between the current branch and master.
func getBranchInfo() []branchInfo {
	info := []branchInfo{}
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
			prExists, prNumber, prState, _ := driver.PullRequestStatus(info[i].Branch, info[i-1].Branch)
			if prExists {
				info[i].PrNumber = prNumber
				info[i].PrState = prState
			}
			branchCommits, _ := script.RunCommandWithCombinedOutput("git", "log", "--no-merges", "--oneline", info[i-1].Branch+".."+info[i].Branch)
			info[i].BranchCommits = strings.Split(strings.TrimSpace(string(branchCommits)), "\n")
		}
	}
	return info
}

func printBranchInfo(branchInfo []branchInfo) {
	var writerFlags uint = 0
	var padChar byte = ' '
	if debugFlag {
		writerFlags |= tabwriter.Debug
		padChar = '.'
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, padChar, writerFlags)
	// fmt.Fprintln(w, "\t Branch\t Github PR\t")

	for i := len(branchInfo) - 1; i >= 0; i-- {
		info := branchInfo[i]
		var branchName, prNumberString string
		if info.IsCurrent {
			branchName = "*" + info.Branch
		} else {
			branchName = info.Branch
		}
		if info.PrNumber != 0 {
			prNumberString = fmt.Sprintf("#%d ", info.PrNumber)
		} else {
			prNumberString = ""
		}
		branchColor := color.New(color.Bold)
		if info.IsCurrent {
			branchColor = branchColor.Add(color.FgCyan)
		}

		// fmt.Fprintln(w, currentMarker, "\t", branchColor.Sprint(info.Branch), "\t ", info.PrInfo, "\t")
		for i, commit := range info.BranchCommits {
			var prefix string
			if i == 0 {
				prefix = fmt.Sprintf("%s\t%s", prNumberString, branchName)
			} else {
				prefix = "\t"
			}
			fmt.Fprintf(w, "%s\t%s\t\n", prefix, commit)
		}
	}
	w.Flush()
}

func init() {
	RootCmd.AddCommand(stackCmd)
}
