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
	path := filepath.Join(self.repository.basePath, self.Name)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	// Stage this change to git repo
	_, err := self.repository.Worktree.Add(".")
	if err != nil {
		return err
	}

	// Commit this change
	err := self.repository.Commit(reason)

	return err
}
