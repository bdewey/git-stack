package steps

import (
	"github.com/bdewey/git-stack/src/script"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoOpStep
}

// Run executes this step.
func (step *DiscardOpenChangesStep) Run() error {
	return script.RunCommand("git", "reset", "--hard")
}
