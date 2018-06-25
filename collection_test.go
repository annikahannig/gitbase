package gitbase

import (
	"os"
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
