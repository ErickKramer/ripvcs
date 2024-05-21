/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <.repos file>",
	Short: "Validate a .repos file",
	Long: `Validate a .repos file.

It checks that all the repositories in the given file have a reachable Git URL
and that the provided version exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Repos file not given!")
			os.Exit(1)
		}
		filePath := args[0]

		// Check that a valid file was given
		config, err := utils.ParseReposFile(filePath)
		if err != nil {
			fmt.Printf("Invalid file given {%s}. %s\n", filePath, err)
			os.Exit(1)
		}

		numWorkers, _ := cmd.Flags().GetInt("workers")
		// Create a channel to send work to the workers with a buffer size of length gitRepos
		jobs := make(chan utils.RepositoryJob, len(config.Repositories))
		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		validFile := false
		for i := 0; i < numWorkers; i++ {
			go func() {
				for job := range jobs {
					validFile = utils.PrintCheckGit(job.DirName, job.Repo.URL, job.Repo.Version, false)
				}
				done <- true
			}()
		}

		for key, repo := range config.Repositories {
			jobs <- utils.RepositoryJob{DirName: key, Repo: repo}
		}
		close(jobs)
		// wait for all goroutines to finish
		for i := 0; i < numWorkers; i++ {
			<-done
		}

		if !validFile {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
}
