package git

import "github.com/bdewey/git-stack/src/command"

// GetLastCommitMessage returns the commit message for the last commit
func GetLastCommitMessage() string {
	return command.New("git", "log", "-1", "--format=%B").Output()
}
