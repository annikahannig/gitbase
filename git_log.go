package gitbase

/*
Go-git supplementals, uses the commandline git interface
and parses the output.

This implements:

  git log --follow <path>

and

  git show <rev>:<path>

*/

import (
	_ "log"
	"os/exec"
	"time"
)

type Commit struct {
	Tree   string
	Parent string

	CreatedAt time.Time
}

func execGitLogFollow(repoPath string, path string) ([]byte, err) {
	cmd := exec.Command(
		"git", "-C", repoPath, "log", "--pretty=raw", "--follow", path,
	)
	return cmd.Output()
}

/*
Parse command output, interpret error
*/
func parseGitLog(data []byte, err error) ([]*Commit, error) {
	return nil, err
}
