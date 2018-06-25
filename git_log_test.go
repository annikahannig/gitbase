package gitbase

import (
	"os"
	"testing"
)

func TestParseGitLog(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	// Add a document
	err = repo.Put("test.doc", []byte("foooo"), "added test document")
	if err != nil {
		t.Error(err)
		return
	}

	// Exec git log
	commits, err := parseGitLog(execGitLogFollow(repo.BasePath, "."))
	if err != nil {
		t.Error(err)
	}

	t.Log(commits)
}

func TestParseTimestampFromAuthor(t *testing.T) {

	// Timestamp: Mon Jun 25 10:46:52 2018 +0200
	author := "Matthias Hannig <matthias@hannig.cc> 1529916412 +0200"

	// Expected UTC Timestamp:
	// 08:46:52 2018 UTC
	createdAt, err := gitParseTimestampFromAuthor(author)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Timestamp:", createdAt)

	if createdAt.Hour() != 8 &&
		createdAt.Minute() != 46 &&
		createdAt.Second() != 52 {
		t.Error("Expected a different result:", createdAt)
	}

}
