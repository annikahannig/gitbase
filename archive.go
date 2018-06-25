package gitbase

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ErrArchiveDoesNotExist = errors.New("archive does not exist")
)

/*
 An Archive represents a sequence of document collections.
 The sequence is derived from the folder name.
*/
type Archive struct {
	Id         uint64
	Collection *Collection
}

func (self *Archive) Path() string {
	if self.Collection == nil {
		return string(self.Id)
	}
	path := filepath.Join(self.Collection.Path(), string(self.Id))
	return path
}

func ArchivePath(collection *Collection, id uint64) string {
	path := filepath.Join(collection.Path(), string(id))
	return path
}

func OpenArchive(collection *Collection, id uint64) (*Archive, error) {
	path := collection.Path()

	// Try to open path
	archivePath := filepath.Join(path, string(id))

	fh, err := os.Open(archivePath)
	if err != nil {
		return nil, ErrArchiveDoesNotExist
	}
	defer fh.Close()

	archive := &Archive{
		Id:         id,
		Collection: collection,
	}

	return archive, nil
}

/*
Calculate next archive id
*/
func NextArchiveId(collection *Collection) uint64 {
	archives, err := ListArchives(collection)
	if err != nil {
		log.Println(err)
		return 1
	}

	maxId := uint64(0)
	for _, archive := range archives {
		if archive.Id > maxId {
			maxId = archive.Id
		}
	}

	return maxId + 1
}

/*
List Archives
*/
func ListArchives(collection *Collection) ([]*Archive, error) {
	archives := []*Archive{}
	path := collection.Path()

	f, err := os.Open(path)
	if err != nil {
		return archives, err
	}
	defer f.Close()

	items, err := f.Readdir(0)
	if err != nil {
		return archives, err
	}

	for _, item := range items {
		if item.IsDir() == false {
			continue
		}

		archiveId, err := strconv.ParseUint(item.Name(), 10, 64)
		if err != nil {
			log.Println("Found non numeric entry in archives path.")
			log.Println("Please check if the repository is OK.")
			continue
		}

		archive := &Archive{
			Id:         archiveId,
			Collection: collection,
		}

		archives = append(archives, archive)
	}

	return archives, nil
}

/*
 List documents
*/
func (self *Archive) Documents() ([]string, error) {
	documents := []string{}

	path := filepath.Join(
		self.Collection.Path(),
		string(self.Id),
	)

	f, err := os.Open(path)
	if err != nil {
		return documents, err
	}
	defer f.Close()

	items, err := f.Readdir(0)
	if err != nil {
		return documents, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		documents = append(documents, item.Name())
	}

	return documents, nil
}

/*
 Remove archive
*/
func (self *Archive) Destroy(reason string) error {
	log.Println("Destroying collection:", self.Name)

	// Fall back to default reason if required
	if reason == "" {
		reason = "removed " + self.Name
	}

	path := filepath.Join(
		self.Collection.Path(),
		string(id),
	)

	fh, err := os.Open(path)
	if err != nil {
		return ErrArchiveDoesNotExist
	}
	defer fh.Close()

	// Disallow write access to repository
	self.Collection.Repository.Lock()
	defer self.Collection.Repository.Unlock()

	// Remove from filesystem
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	// Commit this change
	err = self.Repository.CommitAll(reason)
	return err
}

/*
 Create a new archive with a new id
*/
func CreateArchive(collection *Collection, reason string) (*Archive, error) {
	nextId := NextArchiveId(collection)
	path := ArchivePath(collection, nextId)

	collection.Repository.Lock()
	defer collection.Repository.Unlock()

	// Create if not exists
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	gitkeep := filepath.Join(path, ".gitkeep")
	err = ioutil.WriteFile(gitkeep, []byte{}, 0644)
	if err != nil {
		return nil, err
	}

	err = collection.Repository.CommitAll(reason)
	if err != nil {
		return nil, err
	}

	return OpenArchive(collection, nextId)
}

//
// Wrap document functions
//

/*
 Create / Update document, see Repository.Put
*/
func (self *Archive) Put(key string, document []byte, reason string) error {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.Put(path, document, reason)
}

/*
 Remove document, see: Repository.Remove
*/
func (self *Archive) Remove(key, reason string) error {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.Remove(path, reason)
}

/*
 Fetch, see Repository.Fetch
*/
func (self *Archive) Fetch(key string) ([]byte, error) {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.Fetch(path)
}

/*
 Fetch revision, see Repository.FetchRevision
*/
func (self *Archive) FetchRevision(key, rev string) ([]byte, error) {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.FetchRevision(path, rev)
}

/*
 Get commit History, see Repository.History
*/
func (self *Archive) History(key string) ([]*Commit, error) {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.History(path)
}

/*
 Get revisions, see Repository.Revisions
*/
func (self *Archive) Revisions(key string) ([]string, error) {
	path := filepath.Join(self.Collection.Name, string(self.Id), key)
	return self.Collection.Repository.Revisions(path)
}
