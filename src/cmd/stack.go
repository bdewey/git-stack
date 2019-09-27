package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// branchInfo contains all of the information needed to display a branch.
type branchInfo struct {
	// Branch is the name of the branch.
	Branch string
	// IsCurrent is true if this is the user's current branch.
	IsCurrent bool
	// PrInfo is the name and number of a PR to merge this branch into its parent, if such a PR exists.
	PrInfo string
	// git log --oneline output of the commits that are unique to this branch
	BranchCommits []string
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
			canMerge, defaultCommitMessage, _ := driver.CanMergePullRequest(info[i].Branch, info[i-1].Branch)
			if canMerge {
				info[i].PrInfo = defaultCommitMessage
			}
			branchCommits, _ := script.RunCommandWithCombinedOutput("git", "log", "--no-merges", "--oneline", info[i-1].Branch+".."+info[i].Branch)
			info[i].BranchCommits = strings.Split(string(branchCommits), "\n")
		}
	}
	return info
}

func printBranchInfo(branchInfo []branchInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, ' ', 0)
	fmt.Fprintln(w, "\t Branch\t Github PR")

	for _, info := range branchInfo {
		var currentMarker string
		if info.IsCurrent {
			currentMarker = "*"
		} else {
			currentMarker = " "
		}
		branchColor := color.New(color.Bold)
		if info.IsCurrent {
			branchColor = branchColor.Add(color.FgCyan)
		}
		fmt.Fprintln(w, currentMarker, "\t", branchColor.Sprint(info.Branch), "\t ", info.PrInfo)
		for _, commit := range info.BranchCommits {
			fmt.Fprintln(w, "\t\t ", commit)
		}
	}
	w.Flush()
}

func init() {
	RootCmd.AddCommand(stackCmd)
}
