package gitbase

import (
	"os/exec"

	"errors"
	"fmt"
)

var (
	ErrInvalidRevisionHash = errors.New("invalid revision hash")
	ErrRevsionNotFound     = errors.New("revision not found in repository")
)

func execGitShow(repoPath, path, revision string) ([]byte, error) {
	show := fmt.Sprintf("%s:%s", revision, path)
	cmd := exec.Command(
		"git", "-C", repoPath, "show", show,
	)

	return cmd.Output()
}

// Export
func GitShow(repoPath, path, revision string) ([]byte, error) {
	if !parseGitIsHash(revision) {
		return nil, ErrInvalidRevisionHash
	}
	return execGitShow(repoPath, path, revision)
}
