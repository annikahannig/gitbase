package gitbase

import (
	"errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	_ "io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

/*
A gitbase repository consists of

  * a Repository

    The repository holds a number of Collections


  * Collections
    which will be mapped to named subdirectories

    Example:
    Colection("programs") will be mapped onto

    /path/to/repo/programs


   * Archives
     which essentially are collections of documents, identified by
     a unique (sequential) id

     Example:

     programs, err := repo.Collection("programs")

     program, err := programs.Get(2342)

     An Archive may contain documents:

     source, err := program.Get("source.lua", "HEAD") // []bytes, error

     To create a new version use archive.Put("source.lua", []bytes(content))
     To delete a document, use archive.Delete(key)

*/

var (
	ErrRepositoryPathNotEmpty = errors.New("repository path not empty")
)

type Repository struct {
	sync.RWMutex

	BasePath string
	Worktree *git.Worktree

	gitRepo *git.Repository
}

/*
 Check if the path exists and is empty
*/
func repositoryCanInitialize(path string) error {

	// Check if path exists and is empty
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	items, err := f.Readdir(0)
	if err != nil {
		return err
	}

	if len(items) != 0 {
		return ErrRepositoryPathNotEmpty
	}

	return nil
}

/*
 Open and (if needed) initialize repository
*/
func NewRepository(path string) (*Repository, error) {

	// Assert path exists
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	// Check if we can open this repository
	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		log.Println("Initializing repository:", path)
		err = repositoryCanInitialize(path)
		if err != nil {
			// Path exists, but we can not initialize
			return nil, err
		}

		// Initialize git repo
		gitRepo, err = git.PlainInit(path, false)
		if err != nil {
			return nil, err
		}
	}

	// Open worktree
	worktree, err := gitRepo.Worktree()
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		BasePath: path,
		Worktree: worktree,
		gitRepo:  gitRepo,
	}

	return repo, nil
}

/*
 Stage changes in repository
*/
func (self *Repository) StageChanges() error {
	_, err := self.Worktree.Add(".")
	return err
}

/*
 Commit a change in the repository
*/
func (self *Repository) Commit(reason string) error {
	_, err := self.Worktree.Commit(
		reason, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "gitbase",
				Email: "git@gitbase",
				When:  time.Now(),
			},
		})
	return err
}

/*
 Combined Add + Commit for convenience
*/
func (self *Repository) CommitAll(reason string) error {
	if err := self.StageChanges(); err != nil {
		return err
	}

	return self.Commit(reason)
}

/*
 Get all collections in the repository
*/
func (self *Repository) Collections() []*Collection {

	return nil
}

func (self *Repository) Create(
	name string, reason string,
) (*Collection, error) {
	return CreateCollection(self, name, reason)
}

func (self *Repository) Open(name string) (*Collection, error) {
	return OpenCollection(self, name)
}

func (self *Repository) Use(name string) (*Collection, error) {

	// Try to open collection, if that fails
	collection, err := self.Open(name)
	if err == ErrCollectionDoesNotExist {
		// Try to create the collection
		collection, err = self.Create(
			name, "automatically created collection on use",
		)

		if err != nil {
			return nil, err
		}
	}

	return collection, nil
}
