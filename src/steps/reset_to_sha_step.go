package steps

import (
	"github.com/bdewey/git-stack/src/git"
	"github.com/bdewey/git-stack/src/script"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

// Run executes this step.
func (step *ResetToShaStep) Run() error {
	if step.Sha == git.GetCurrentSha() {
		return nil
	}
	cmd := []string{"git", "reset"}
	if step.Hard {
		cmd = append(cmd, "--hard")
	}
	cmd = append(cmd, step.Sha)
	return script.RunCommand(cmd...)
}
