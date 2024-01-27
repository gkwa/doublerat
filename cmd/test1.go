/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/doublerat/workflow1"
)

// test1Cmd represents the test1 command
var test1Cmd = &cobra.Command{
	Use:   "test1",
	Short: "Initialize a new Git repository with a submodule",
	Long: `This command initializes a new Git repository at /tmp/test and adds a submodule
(https://github.com/taylormonacelli/darksheep) to the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := workflow1.RunTest()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	// Define flags and configuration settings here.
}
