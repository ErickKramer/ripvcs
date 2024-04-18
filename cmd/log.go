/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Get logs of found repositories.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var root string
		if len(args) == 0 {
			root = "."
		} else {
			root = args[0]
		}
		gitRepos := utils.FindGitRepositories(root)

		onelineFlag, _ := cmd.Flags().GetBool("oneline")
		numWorkers, _ := cmd.Flags().GetInt("workers")
		numCommits, _ := cmd.Flags().GetInt("num-commits")

		// Create a channel to send work to the workers with a buffer size of length gitRepos
		// HINT: The buffer size specifies how many elements the channel can hold before blocking sends
		jobs := make(chan string, len(gitRepos))

		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		// Iterate over the numWorkers
		for i := 0; i < numWorkers; i++ {
			go func() {
				for repo := range jobs {
					utils.PrintGitLog(repo, onelineFlag, numCommits)
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
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().IntP("workers", "w", 8, "Number of workers to use for concurrency")
	logCmd.Flags().IntP("num-commits", "n", 4, "Show only the last n commits")
	logCmd.Flags().BoolP("oneline", "l", false, "Show short version of logs")
}
