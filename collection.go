package gitbase

import (
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
	Name string

	repository *Repository
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

	// Disallow write access to repository
	self.repository.Lock()
	defer self.repository.Unlock()

	// Remove from filesystem
	err := os.RemoveAll(self.Path())
	if err != nil {
		return err
	}

	// Stage this change to git repo
	if err = self.repository.StageChanges(); err != nil {
		return err
	}

	// Commit this change
	err := self.repository.Commit(reason)

	return err
}

/*
 Calculate path of collection, derived from
 Name and the collection's base path
*/
func (self *Collection) Path() string {
	basePath := ""
	if self.repository != nil {
		basePath = self.repository.BasePath
	}

	return filepath.Join(basePath, self.Name)
}
