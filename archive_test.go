package gitbase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNextArchiveId(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	foo, err := repo.Use("foo")

	// Next id should be 1 (empty, and first insert)
	nextId := NextArchiveId(foo)
	if nextId != 1 {
		t.Error("Next id should be 1")
		return
	}

	// Create a bogus archive
	archivePath := filepath.Join(foo.Path(), "22")
	os.MkdirAll(archivePath, 0755)

	// Next id should be 22 + 1
	nextId = NextArchiveId(foo)
	if nextId != 23 {
		t.Error("Expected NextId to be 23, got:", nextId)
		return
	}

}
