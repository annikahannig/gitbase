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

	err = repo.Put("test2.doc", []byte("bar"), "added\nanother\ntest document")
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

	if len(commits) != 2 {
		t.Error("Expected 2 commits!")
	}

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

func TestParseGitIsHash(t *testing.T) {
	tests := map[string]bool{
		"d7585cbdf989fa9ddd810aeb08ee41c11fbca8bb": true,
		"7da3af6390f7a400c6265f98768ed595bb477b8b": true,
		"400c6265f98768ed5":                        true,
		"1234567890abcdef":                         true,
		"fail":                                     false,
		"2342!":                                    false,
	}

	for hash, expected := range tests {
		if result := parseGitIsHash(hash); result != expected {
			t.Error("Expected:", expected, "got:", result)
		}
	}
}
