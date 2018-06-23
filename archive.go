package gitbase

import (
	"errors"
	_ "io/ioutil"
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

func FindArchive(collection *Collection, id uint64) (*Archive, error) {
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
Create a new archive with a new id
*/
func CreateArchive(collection *Collection, reason string) (*Archive, error) {
	nextId := collection.NextId()

	_ = nextId

	return nil, nil
}
