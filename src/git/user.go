package git

import "github.com/bdewey/git-stack/src/command"

// GetLocalAuthor returns the locally Git configured user
func GetLocalAuthor() string {
	name := command.New("git", "config", "user.name").Output()
	email := command.New("git", "config", "user.email").Output()
	return name + " <" + email + ">"
}
