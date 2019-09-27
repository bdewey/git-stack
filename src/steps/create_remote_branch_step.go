package steps

import (
	"github.com/bdewey/git-stack/src/script"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	NoOpStep
	BranchName string
	Sha        string
}

// Run executes this step.
func (step *CreateRemoteBranchStep) Run() error {
	return script.RunCommand("git", "push", "origin", step.Sha+":refs/heads/"+step.BranchName)
}
