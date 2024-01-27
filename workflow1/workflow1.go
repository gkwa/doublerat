package workflow1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

const repoDir = "mytest1706299472"

const gitmodTemplate = `[submodule "{{.Name}}"]
	path = {{.Path}}
	url = {{.URL}}
	active = true
`

func RunTest() error {
	fmt.Printf("Deleting directory: %s\n", repoDir)
	if err := os.RemoveAll(repoDir); err != nil {
		return fmt.Errorf("error deleting directory: %v", err)
	}

	if _, err := os.Stat(repoDir); err == nil {
		return fmt.Errorf("assertion failed: directory %s still exists after removal", repoDir)
	}

	fmt.Printf("Directory %s successfully removed\n", repoDir)

	fmt.Printf("Creating Git repository: %s\n", repoDir)
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		return fmt.Errorf("error initializing Git repository: %v", err)
	}
	fmt.Println("Git repository created successfully.")

	var repoService RepositoryService
	var repos []RepositoryInfo

	repoService = &StaticRepositoryService{}
	repos, err = retrieveRepositories(repoService)
	if err != nil {
		return fmt.Errorf("error retrieving repositories using static strategy: %v", err)
	}

	// repoService = &JSONFileRepositoryService{
	// 	FilePath: "/Users/mtm/pdev/taylormonacelli/hisrabbit/data1.json",
	// }
	// repos, err = retrieveRepositories(repoService)
	// if err != nil {
	// 	return fmt.Errorf("error retrieving repositories using JSON file strategy: %v", err)
	// }

	for _, r := range repos {
		fmt.Printf("Adding submodule: %s\n", r.Path)
		err = AddSubmodule(repo, r.Path, r.Path, r.GitURL, r.Version)
		if err != nil {
			return fmt.Errorf("error adding submodule: %v", err)
		}
	}

	return nil
}

// https://github.com/go-git/go-git/issues/212#issuecomment-757537285
// AddSubmodule adds a new git submodule.
func AddSubmodule(repo *git.Repository, name, path, url, branch string) error {
	spec := config.Submodule{
		Name:   name,
		Path:   path,
		URL:    url,
		Branch: branch,
	}

	slog.Debug("Adding submodule", "name", name, "path", path, "url", url, "branch", branch)

	wtree, err := repo.Worktree()
	if err != nil {
		return err
	}

	gitmodulesFile := filepath.Join(wtree.Filesystem.Root(), ".gitmodules")
	f, err := os.OpenFile(gitmodulesFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	t := template.Must(template.New("gitmodule").Parse(gitmodTemplate))
	if err := t.Execute(f, spec); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	submod, err := wtree.Submodule(name)
	if err != nil {
		return err
	}

	if err := submod.Init(); err != nil {
		return err
	}

	subrepo, err := submod.Repository()
	if err != nil {
		return err
	}

	subwtree, err := subrepo.Worktree()
	if err != nil {
		return err
	}

	opts := &git.PullOptions{RemoteName: "origin"}
	if err := subwtree.Pull(opts); err != nil {
		return err
	}

	return nil
}

// RepositoryService defines the methods for interacting with repositories.
type RepositoryService interface {
	GetRepositories() ([]RepositoryInfo, error)
}

// StaticRepositoryService is a RepositoryService implementation that uses a static list.
type StaticRepositoryService struct{}

// GetRepositories retrieves repositories from a static list.
func (s *StaticRepositoryService) GetRepositories() ([]RepositoryInfo, error) {
	repositoryInfoSlice := []RepositoryInfo{
		{
			BrowseURL: "https://github.com/shykes/daggerverse/tree/792b8449b95393100866720e26a27e34d818738a/ollama",
			CreatedAt: time.Date(2024, 0o1, 25, 8, 52, 15, 0, time.UTC),
			GitCommit: "792b8449b95393100866720e26a27e34d818738a",
			GitURL:    "https://github.com/shykes/daggerverse",
			IndexedAt: time.Date(2024, 0o1, 25, 8, 52, 34, 724422000, time.UTC),
			Path:      "github.com/shykes/daggerverse/ollama",
			Release:   "",
			Subpath:   "ollama",
			Version:   "792b8449b95393100866720e26a27e34d818738a",
		},
		{
			BrowseURL: "https://github.com/samalba/inline-python-mod",
			CreatedAt: time.Date(2024, 0o1, 25, 3, 4, 50, 0, time.UTC),
			GitCommit: "e0932748103867e73f6d63165823d4a830cd358c",
			GitURL:    "https://github.com/samalba/inline-python-mod",
			IndexedAt: time.Date(2024, 0o1, 25, 3, 8, 33, 26385000, time.UTC),
			Path:      "github.com/samalba/inline-python-mod",
			Release:   "",
			Version:   "e0932748103867e73f6d63165823d4a830cd358c",
		},
		{
			BrowseURL: "https://github.com/samalba/inline-python-mod",
			CreatedAt: time.Date(2024, 0o1, 25, 2, 58, 31, 0, time.UTC),
			GitCommit: "2f5521e9bd43e5ac3219a36f930c0581dc299c74",
			GitURL:    "https://github.com/samalba/inline-python-mod",
			IndexedAt: time.Date(2024, 0o1, 25, 2, 58, 59, 382617000, time.UTC),
			Path:      "github.com/samalba/inline-python-mod",
			Release:   "",
			Version:   "2f5521e9bd43e5ac3219a36f930c0581dc299c74",
		},
	}

	return repositoryInfoSlice, nil
}

// JSONFileRepositoryService is a RepositoryService implementation that reads data from a JSON file.
type JSONFileRepositoryService struct {
	FilePath string
}

// GetRepositories retrieves repositories from a JSON file.
func (j *JSONFileRepositoryService) GetRepositories() ([]RepositoryInfo, error) {
	file, err := os.Open(j.FilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %v", err)
	}
	defer file.Close()

	var repositoryInfoSlice []RepositoryInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&repositoryInfoSlice)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON file: %v", err)
	}

	return repositoryInfoSlice, nil
}

func retrieveRepositories(repoService RepositoryService) ([]RepositoryInfo, error) {
	return repoService.GetRepositories()
}
