package steps

import "github.com/bdewey/git-stack/src/script"

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortMergeBranchStep) Run() error {
	return script.RunCommand("git", "merge", "--abort")
}
