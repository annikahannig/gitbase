package gitbase

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testRepoPath() string {
	return filepath.Join(os.TempDir(), "gitbase-test-repo")
}

func TestRepositoryCanInitialize(t *testing.T) {
	path := testRepoPath()

	// This should fail:
	err := repositoryCanInitialize(path)
	if err == nil {
		t.Error("The path should not be initializable")
		return
	} else {
		t.Log("Initialization failed (GOOD!), reason:", err)
	}

	// Let's create this path and try again
	err = os.MkdirAll(path, 0755)
	if err != nil {
		t.Error("Could not create test repo path:", err)
		return
	}
	defer os.RemoveAll(path)

	// This should not fail:
	err = repositoryCanInitialize(path)
	if err != nil {
		t.Error("The path should be initializable!")
		t.Error(err)
		return
	}

	// Let's make this path not empty
	if err = ioutil.WriteFile(
		filepath.Join(path, "fail_file"),
		[]byte{}, 0644); err != nil {

		t.Error(err)
		return
	}

	// This should fail with Path Not Empty error
	err = repositoryCanInitialize(path)
	if err != ErrRepositoryPathNotEmpty {
		t.Error("Expected ErrRepositoryPathNotEmpty, got:", err)
		return
	}
}

func TestRepositoryInitialization(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	_, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	// This should work aswell, because the repo should already
	// be initialized
	_, err = NewRepository(path)
	if err != nil {
		t.Error("Could not open repo:", err)
		return
	}
}

func TestRepositoryDocumentStorage(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path)

	repo, err := NewRepository(path)
	if err != nil {
		t.Error(err)
		return
	}

	// Okay try to add a simple document
	document := []byte("Hello World!")

	// This should fail, because it should not exist
	_, err = repo.Fetch("hello.doc")
	if err == nil {
		t.Error("Expected hello.doc to not exit!")
	}

	// Create our test document
	err = repo.Put("hello.doc", document, "added test document")
	if err != nil {
		t.Error(err)
	}

	// Fetch this
	retrieved, err := repo.Fetch("hello.doc")
	if err != nil {
		t.Error(err)
	}

	if string(retrieved) != string(document) {
		t.Error("Retrieved document differs from added document")
	}

	// Remove document
	err = repo.Remove("hello.doc", "not longer required")
	if err != nil {
		t.Error(err)
	}

	// This should fail
	_, err = repo.Fetch("hello.doc")
	if err == nil {
		t.Error("Expected fetch(hello.doc) to fail after removal!")
	}

}
