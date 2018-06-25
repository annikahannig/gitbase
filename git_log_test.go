package gitbase

import (
	"os"
	"testing"
)

func TestExecGitLog(t *testing.T) {
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
	execGitLog(repo.BasePath, ".")

}
