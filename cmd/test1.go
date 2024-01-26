/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"

	"github.com/spf13/cobra"
)

// test1Cmd represents the test1 command
var test1Cmd = &cobra.Command{
	Use:   "test1",
	Short: "Initialize a new Git repository with a submodule",
	Long: `This command initializes a new Git repository at /tmp/test and adds a submodule
(https://github.com/taylormonacelli/darksheep) to the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTest1()
	},
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	// Define flags and configuration settings here.
}

func runTest1() error {
	// Initialize a new git repository
	repo, err := git.PlainClone(".", false, &git.CloneOptions{
		URL:               "https://github.com/taylormonacelli/darksheep",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	// Get the submodule
	submoduleName := "darksheep"

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}

	sub, err := w.Submodule(submoduleName)
	if err != nil {
		return fmt.Errorf("error getting submodule: %v", err)
	}

	sr, err := sub.Repository()
	if err != nil {
		return fmt.Errorf("error getting submodule repository: %v", err)
	}

	sw, err := sr.Worktree()
	if err != nil {
		return fmt.Errorf("error getting submodule worktree: %v", err)
	}

	// Pull the latest changes in the submodule
	fmt.Println("git submodule update --remote")
	err = sw.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		return fmt.Errorf("error pulling submodule changes: %v", err)
	}

	fmt.Println("Git repository with submodule created successfully.")

	return nil
}
