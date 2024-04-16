/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of any relative repository found relative to the given path.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1), // Specify positional argument for the path
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		if len(args) == 0 {
			root = "."
		} else {
			root = args[0]
		}
		gitRepos := utils.FindGitRepositories(root)

		skipEmtpy, _ := cmd.Flags().GetBool("skip-empty")
		numWorkers, _ := cmd.Flags().GetInt("workers")

		// Create a channel to send work to the workers with a buffer size of length gitRepos
		// HINT: The buffer size specifies how many elements the channel can hold before blocking sends
		jobs := make(chan string, len(gitRepos))

		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		// Iterate over the numWorkers
		for i := 0; i < numWorkers; i++ {
			go func() {
				for repo := range jobs {
					utils.PrintGitStatus(repo, skipEmtpy)
				}
				done <- true
			}()
		}
		// Send each git repository path to the jobs channel
		for _, repo := range gitRepos {
			jobs <- repo
		}
		close(jobs) // Close channel to signal no more work will be sent

		// wait for all goroutines to finish
		for i := 0; i < numWorkers; i++ {
			<-done
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	statusCmd.Flags().IntP("workers", "w", 8, "Number of workers to use for concurrency")
	statusCmd.Flags().BoolP("skip-empty", "s", false, "Skip repositories with clean working tree.")
	// TODO: Add workers functionality to handle the repositories with go routines
}
