package steps

import (
	"github.com/bdewey/git-stack/src/drivers"
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/script"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *CreatePullRequestStep) Run() error {
	driver := drivers.GetActiveDriver()
	parentBranch := git.GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestURL(step.BranchName, parentBranch))
	return nil
}
