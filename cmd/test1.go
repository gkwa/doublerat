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
		// runTest1()
		runTest2()
	},
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	// Define flags and configuration settings here.
}

const repoDir = "mytest1706299472"

func runTest1() {
	fmt.Printf("Deleting directory: %s\n", repoDir)
	if err := os.RemoveAll(repoDir); err != nil {
		fmt.Printf("Error deleting directory: %v\n", err)
		return
	}

	fmt.Printf("Creating Git repository: %s\n", repoDir)
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		fmt.Printf("Error initializing Git repository: %v\n", err)
		return
	}
	fmt.Println("Git repository created successfully.")

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Error getting worktree: %v\n", err)
		return
	}

	sm, err := wt.Submodule("basic")
	if err != nil {
		fmt.Printf("Error getting submodule: %v\n", err)
		return
	}

	name := "darksheep"
	path := "darksheep"
	url := "https://github.com/taylormonacelli/darksheep"
	branch := "master"

	submod := &config.Submodule{
		Name:   name,
		Path:   path,
		URL:    url,
		Branch: branch,
	}

	rConfig, err := repo.Config()
	if err != nil {
		fmt.Println("Error getting repo config")
	}

	rConfig.Submodules[name] = submod

	if err := repo.SetConfig(rConfig); err != nil {
		fmt.Println("Error setting repo config")
		return
	}

	if err := sm.Init(); err != nil {
		fmt.Printf("Error initializing submodule: %v\n", err)
		return
	}

	subrepo, err := sm.Repository()
	if err != nil {
		fmt.Printf("Error getting submodule repository: %v\n", err)
		return
	}

	subwtree, err := subrepo.Worktree()
	if err != nil {
		fmt.Printf("Error getting submodule worktree: %v\n", err)
		return
	}

	opts := &git.PullOptions{RemoteName: "origin"}
	if err := subwtree.Pull(opts); err != nil {
		fmt.Printf("Error pulling submodule: %v\n", err)
		return
	}

	fmt.Println("Submodule added successfully.")
}

const gitmodTemplate = `[submodule "{{.Name}}"]
	path = {{.Path}}
	url = {{.URL}}
	active = true
`

func runTest2() {
	fmt.Printf("Deleting directory: %s\n", repoDir)
	if err := os.RemoveAll(repoDir); err != nil {
		fmt.Printf("Error deleting directory: %v\n", err)
		return
	}

	fmt.Printf("Creating Git repository: %s\n", repoDir)
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		fmt.Printf("Error initializing Git repository: %v\n", err)
		return
	}
	fmt.Println("Git repository created successfully.")

	name := "darksheep"
	path := "darksheep"
	url := "https://github.com/taylormonacelli/darksheep"
	branch := "master"

	err = AddSubmodule(repo, name, path, url, branch)
	if err != nil {
		fmt.Printf("Error adding submodule: %v\n", err)
		return
	}
}

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
	f, err := os.OpenFile(gitmodulesFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
