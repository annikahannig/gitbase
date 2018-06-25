package gitbase

import (
	"os"
	"testing"
)

func TestExecGitShow(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	// Add a document
	err = repo.Put("test.doc", []byte("fo\nooo"), "added test document")
	if err != nil {
		t.Error(err)
		return
	}

	err = repo.Put("test.doc", []byte("bar"), "updated test document")
	if err != nil {
		t.Error(err)
		return
	}

	// There should be two revisions for test.doc
	revisions, err := repo.Revisions("test.doc")
	if err != nil {
		t.Error(err)
		return
	}

	// Test git show
	result, err := execGitShow(repo.BasePath, "test.doc", revisions[0])
	if err != nil {
		t.Error(err)
	}

	if string(result) != "bar" {
		t.Error("Expected: 'bar', got:", string(result))
	}

	result, err = execGitShow(repo.BasePath, "test.doc", revisions[1])
	if err != nil {
		t.Error(err)
	}
	if string(result) != "fo\nooo" {
		t.Error("Expected: 'fo\\nooo', got:", string(result))
	}

	// This should yield an error
	result, err = execGitShow(repo.BasePath, "test.doc", "d3adb33f")
	if err == nil {
		t.Error(err)
	}

	_, err = execGitShow(repo.BasePath, "test.dog", revisions[0])
	if err == nil {
		t.Error("Expected error with unkonwn file")
	}

}
