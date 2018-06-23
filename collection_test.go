package gitbase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectionCreateDestroy(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	// This should fail
	_, err = OpenCollection(repo, "test")
	if err != ErrCollectionDoesNotExist {
		t.Error("Expected:", ErrCollectionDoesNotExist, "got:", err)
	}

	// This should work
	collection, err := CreateCollection(repo, "test", "new test collection")
	if err != nil {
		t.Error(err)
		return
	}

	// This work aswell
	testCollection, err := CreateCollection(repo, "test", "whatever")
	if err != nil {
		t.Error(err)
		return
	}

	// Remove collection
	err = collection.Destroy("test reason")
	if err != nil {
		t.Error(err)
	}

	// This should fail however:
	err = testCollection.Destroy("another reason")
	if err == nil {
		t.Error("A collection should not be removable twice")
	}

}

// Id Sequence

func TestCollectionNextId(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	foo, err := repo.Use("foo")

	// Next id should be 1 (empty, and first insert)
	nextId := foo.NextId()
	if nextId != 1 {
		t.Error("Next id should be 1")
		return
	}

	// Create a bogus archive
	archivePath := filepath.Join(foo.Path(), "22")
	os.MkdirAll(archivePath, 0755)

	// Next id should be 22 + 1
	nextId = foo.NextId()
	if nextId != 23 {
		t.Error("Expected NextId to be 23, got:", nextId)
		return
	}

}
