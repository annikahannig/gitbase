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
	"os/exec"

	"bufio"
	"bytes"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidAuthorTimestamp = errors.New("invalid timestamp in line")
)

type Commit struct {
	Id     string
	Tree   string
	Parent string

	Author    string
	Committer string

	Message string

	CreatedAt time.Time
}

func execGitLogFollow(repoPath string, path string) ([]byte, error) {
	cmd := exec.Command(
		"git", "-C", repoPath, "log", "--pretty=raw", "--follow", path,
	)
	return cmd.Output()
}

/*
Identify new commit
*/
func parseGitIsHeaderStart(line string) bool {
	tokens := strings.Split(line, " ")
	if len(tokens) != 2 {
		return false
	}

	if tokens[0] != "commit" {
		return false
	}

	return parseGitIsHash(tokens[1])
}

func parseGitIsHash(hash string) bool {
	sigma := "0123456789abcdef"
	for _, c := range hash {
		ok := false
		for _, t := range sigma {
			if c == t {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}

	return true
}

/*
Parse command output, interpret error
*/
func parseGitLog(data []byte, err error) ([]*Commit, error) {
	commits := []*Commit{}
	if err != nil {
		return commits, err
	}

	// States:
	stateHeader := 1
	stateMessage := 2
	state := stateHeader

	// Current commit
	var commit *Commit

	// Process data linewise
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			log.Println(err)
			continue
		}

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if parseGitIsHeaderStart(line) {
			if commit == nil {
				// First commit
				commit = &Commit{}
			} else {
				// Next commit
				commit.Message = strings.TrimSpace(commit.Message)
				commits = append(commits, commit)
				commit = &Commit{}
			}

			state = stateHeader
		}

		if state == stateHeader {
			tokens := strings.SplitN(line, " ", 2)
			switch tokens[0] {
			case "":
				// Separator, next state: message
				state = stateMessage

				// Set created at time by parsing the
				// author line
				createdAt, err := gitParseTimestampFromAuthor(commit.Author)
				if err != nil {
					log.Println(err)
					continue
				}
				commit.CreatedAt = createdAt

			case "commit":
				commit.Id = tokens[1]
				break
			case "tree":
				commit.Tree = tokens[1]
				break
			case "parent":
				commit.Parent = tokens[1]
				break
			case "author":
				commit.Author = tokens[1]
				break
			case "committer":
				commit.Committer = tokens[1]
				break
			default:
				log.Println("Unknown token:", tokens[0])
			}
		} else if state == stateMessage {
			commit.Message += line + "\n"
		}
	}

	// Add last commit
	if commit != nil {
		commit.Message = strings.TrimSpace(commit.Message)
		commits = append(commits, commit)
	}

	return commits, nil
}

func gitParseTimestampFromAuthor(line string) (time.Time, error) {
	tokens := strings.Split(line, " ")
	tlen := len(tokens)
	if tlen < 2 {
		return time.Unix(0, 0), ErrInvalidAuthorTimestamp
	}

	offset := tokens[tlen-1]
	timestamp := tokens[tlen-2]

	// Parse offset
	loc, err := DatetimeParseOffset(offset)
	if err != nil {
		return time.Unix(0, 0), err
	}

	// Parse timestamp
	timestampUnix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}

	createdAt := DatetimeSetLocation(
		time.Unix(timestampUnix, 0),
		loc,
	).UTC()

	return createdAt, nil
}

func GitHistory(basePath, path string) ([]*Commit, error) {
	return parseGitLog(execGitLogFollow(basePath, path))
}
