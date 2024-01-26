/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

// test1Cmd represents the test1 command
var test1Cmd = &cobra.Command{
	Use:   "test1",
	Short: "Initialize a new Git repository with a submodule",
	Long: `This command initializes a new Git repository at /tmp/test and adds a submodule
(https://github.com/taylormonacelli/darksheep) to the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTest()
	},
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	// Define flags and configuration settings here.
}

const repoDir = "mytest1706299472"

const gitmodTemplate = `[submodule "{{.Name}}"]
	path = {{.Path}}
	url = {{.URL}}
	active = true
`

func runTest() {
	fmt.Printf("Deleting directory: %s\n", repoDir)
	if err := os.RemoveAll(repoDir); err != nil {
		fmt.Printf("Error deleting directory: %v\n", err)
		return
	}

	if _, err := os.Stat(repoDir); err == nil {
		fmt.Printf("Assertion failed: Directory %s still exists after removal\n", repoDir)
		return
	}

	fmt.Printf("Directory %s successfully removed\n", repoDir)

	fmt.Printf("Creating Git repository: %s\n", repoDir)
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		fmt.Printf("Error initializing Git repository: %v\n", err)
		return
	}
	fmt.Println("Git repository created successfully.")

	err = AddSubmodule(repo, "darksheep", "darksheep", "https://github.com/taylormonacelli/darksheep", "master")
	if err != nil {
		fmt.Printf("Error adding submodule: %v\n", err)
		return
	}

	err = AddSubmodule(repo, "greenleeks", "greenleeks", "https://github.com/taylormonacelli/greenleeks", "master")
	if err != nil {
		fmt.Printf("Error adding submodule: %v\n", err)
		return
	}
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
