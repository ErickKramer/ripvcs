/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
					invalidFile = utils.PrintCheckGit(job.DirName, job.Repo.URL, job.Repo.Version, false)
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
	validateCmd.Flags().IntP("workers", "w", 8, "Number of workers to use for concurrency")
}
