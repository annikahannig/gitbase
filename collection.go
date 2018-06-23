package gitbase

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/*
 A collection represents a storage for archives
 this is basically a folder at the repositories base
 path.
*/

type Collection struct {
	Name       string
	Repository *Repository
}

var (
	ErrCollectionDoesNotExist = errors.New("collection does not exist")
)

/*
 Calculate path of collection, derived from
 Name and the collection's base path
*/
func (self *Collection) Path() string {
	basePath := ""
	if self.Repository != nil {
		basePath = self.Repository.BasePath
	}

	return filepath.Join(basePath, self.Name)
}

/*
 Remove collection from repository
*/
func (self *Collection) Destroy(reason string) error {
	log.Println("Destroying collection:", self.Name)

	// Fall back to default reason if required
	if reason == "" {
		reason = "removed " + self.Name
	}

	path := self.Path()

	// Check if we can remove
	fh, err := os.Open(path)
	if err != nil {
		return ErrCollectionDoesNotExist
	}
	defer fh.Close()

	// Disallow write access to repository
	self.Repository.Lock()
	defer self.Repository.Unlock()

	// Remove from filesystem
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	// Stage this change to git repo
	if err = self.Repository.StageChanges(); err != nil {
		return err
	}

	// Commit this change
	err = self.Repository.Commit(reason)

	return err
}

/*
 Create Collection
*/
func CreateCollection(
	repo *Repository,
	name string,
	reason string,
) (*Collection, error) {
	collection := &Collection{
		Name:       name,
		Repository: repo,
	}
	// Lock repository
	repo.Lock()
	defer repo.Unlock()

	// Create filesystem path
	path := collection.Path()
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	// Add .gitkeep (for now, to have something to add
	// to the repo). In future consider creating some
	// metadata document.
	gitkeep := filepath.Join(path, ".gitkeep")
	err = ioutil.WriteFile(gitkeep, []byte{}, 0644)
	if err != nil {
		return nil, err
	}

	// Consider adding document storage support to
	// collections.

	// Insert into repository
	if err = repo.CommitAll(reason); err != nil {
		return nil, err
	}

	return collection, nil
}

/*
 Open collection in repository
*/
func OpenCollection(
	repo *Repository,
	name string,
) (*Collection, error) {
	collection := &Collection{
		Name:       name,
		Repository: repo,
	}
	path := collection.Path()

	// Check if collection exists
	fh, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, ErrCollectionDoesNotExist
	}
	defer fh.Close()

	// Great, file exists, peachy.
	return collection, nil
}

/*
 Get all archives
*/
func (self *Collection) Archives() ([]*Archive, error) {
	return ListArchives(self)
}

/*
 Calculate next id
*/
func (self *Collection) NextId() uint64 {
	return NextArchiveId(self)
}

/*
 Find Archive
*/
func (self *Collection) Find(id uint64) (*Archive, error) {
	return OpenArchive(self, id)
}
