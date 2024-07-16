/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <optional path>",
	Short: "Check status of all repositories",
	Long: `Check status of all repositories.

If no path is given, it checks the status of any Git repository relative to the current path.`,
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		if len(args) == 0 {
			root = "."
		} else {
			root = utils.GetRepoPath(args[0])
		}
		gitRepos := utils.FindGitRepositories(root)

		plainStatus, _ := cmd.Flags().GetBool("plain")
		skipEmtpy, _ := cmd.Flags().GetBool("skip-empty")
		numWorkers, _ := cmd.Flags().GetInt("workers")

		// Create a channel to send work to the workers with a buffer size of length gitRepos
		jobs := make(chan string, len(gitRepos))

		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		// Iterate over the numWorkers
		for i := 0; i < numWorkers; i++ {
			go func() {
				for repo := range jobs {
					utils.PrintGitStatus(repo, skipEmtpy, plainStatus)
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
	statusCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
	statusCmd.Flags().BoolP("plain", "p", false, "Show simpler status report")
	statusCmd.Flags().BoolP("skip-empty", "s", false, "Skip repositories with clean working tree.")
}
