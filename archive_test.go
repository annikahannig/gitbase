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
	if err != nil {
		t.Error(err)
	}

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

func TestArchiveCreateDestroy(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	collection, err := repo.Use("fnord")
	if err != nil {
		t.Error(err)
		return
	}

	archive1, err := collection.NextArchive("make test archive")
	if err != nil {
		t.Error(err)
		return
	}
	archive2, err := collection.NextArchive("make another test archive")
	if err != nil {
		t.Error(err)
		return
	}

	err = archive1.Destroy("not required anymore")
	if err != nil {
		t.Error(err)
	}

	ret, err := collection.Find(2)
	if err != nil {
		t.Error(err)
		return
	}

	if ret.Id != archive2.Id {
		t.Error("Expected retrieved archive id to be 2, got:", ret.Id)
	}

}

func TestArchiveDocumentHandling(t *testing.T) {
	path := testRepoPath()
	defer os.RemoveAll(path) // Clean up afterwards

	repo, err := NewRepository(path)
	if err != nil {
		t.Error("Could not initialize repo:", err)
		return
	}

	collection, err := repo.Use("foo")
	if err != nil {
		t.Error(err)
	}

	archive, err := collection.NextArchive("new test archive")
	if err != nil {
		t.Error(err)
	}

	if archive.Id != 1 {
		t.Error("Expected archive Id: 1")
	}

	// Add documents
	docHello := []byte("hello document")
	docFoo := []byte("foo bar baz")

	err = archive.Put("hello", docFoo, "added test document")
	if err != nil {
		t.Error(err)
	}
	err = archive.Put("hello", docHello, "updated test document")
	if err != nil {
		t.Error(err)
	}
	err = archive.Put("foo", docFoo, "added another test document")
	if err != nil {
		t.Error(err)
	}

	// Try a simple fetch
	res, err := archive.Fetch("foo")
	if err != nil {
		t.Error(err)
	}
	if string(res) != string(docFoo) {
		t.Error(
			"Retrieved document does not match put document:",
			string(res),
		)
	}

	// A simple delete
	err = archive.Remove("foo", "done with this")
	if err != nil {
		t.Error(err)
	}

	// Let's get the history of hello
	history, err := archive.History("hello")
	if len(history) != 2 {
		t.Error("History should have two entries")
	}

	if history[1].Message != "added test document" {
		t.Error(
			"Expected a different commit message, got:",
			history[1].Message,
		)
	}

	revs, err := archive.Revisions("hello")
	if err != nil {
		t.Error(err)
	}

	res, err = archive.FetchRevision("hello", revs[1])
	if err != nil {
		t.Error(err)
	}
	if string(res) != string(docFoo) {
		t.Error("Expected foo, got:", string(res))
	}
}
