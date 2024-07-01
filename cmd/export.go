/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var exportCmd = &cobra.Command{
	Use:   "export <optional path>",
	Short: "Export list of available repositories",
	Long: `Export list of available repositories..

If no path is given, it checks the finds any Git repository relative to the current path.`,
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		if len(args) == 0 {
			root = "."
		} else {
			root = args[0]
		}
		gitRepos := utils.FindGitRepositories(root)

		filePath, _ := cmd.Flags().GetString("output")
		numWorkers, _ := cmd.Flags().GetInt("workers")
		getCommitsFlag, _ := cmd.Flags().GetBool("commits")
		getTagsFlags, _ := cmd.Flags().GetBool("tags")

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
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
	exportCmd.Flags().StringP("output", "o", "", "Path to output `.repos` file")
	exportCmd.Flags().BoolP("commits", "c", false, "Export repositories hashes instead of branches")
	exportCmd.Flags().BoolP("tags", "t", false, "Export repositories tags instead of branches.")
}
