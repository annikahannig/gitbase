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

func TestCanInitialize(t *testing.T) {
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