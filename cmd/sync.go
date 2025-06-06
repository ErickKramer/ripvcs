/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// syncCmd represents the pull command
var syncCmd = &cobra.Command{
	Use:   "sync <optional path>",
	Short: "Synchronize all found repositories.",
	Long: `Synchronize all found repositories.

It stashes all changes found in the repostory, pull latest remote,
and bring back staged changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		if len(args) == 0 {
			root = "."
		} else {
			root = utils.GetRepoPath(args[0])
		}
		gitRepos := utils.FindGitRepositories(root)

		numWorkers, _ := cmd.Flags().GetInt("workers")

		// Create a channel to send work to the workers with a buffer size of length gitRepos
		// HINT: The buffer size specifies how many elements the channel can hold before blocking sends
		jobs := make(chan string, len(gitRepos))

		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		// Iterate over the numWorkers
		for range numWorkers {
			go func() {
				for repo := range jobs {
					utils.PrintGitSync(repo)
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
		for range numWorkers {
			<-done
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
}
