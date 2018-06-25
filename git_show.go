package gitbase

import (
	"os/exec"

	"errors"
	"fmt"
)

var (
	ErrRevsionNotFound = errors.New("revision not found in repository")
)

func execGitShow(repoPath, path, revision string) ([]byte, error) {
	show := fmt.Sprintf("%s:%s", revision, path)
	cmd := exec.Command(
		"git", "-C", repoPath, "show", show,
	)

	return cmd.Output()
}
