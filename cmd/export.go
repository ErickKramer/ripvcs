/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"ripvcs/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
			root = utils.GetRepoPath(args[0])
		}
		gitRepos := utils.FindGitRepositories(root)

		filePath, _ := cmd.Flags().GetString("output")
		visualizeOutput, _ := cmd.Flags().GetBool("visualize")

		skipOutputFile := false

		if len(filePath) == 0 {
			if visualizeOutput {
				skipOutputFile = true
			} else {
				utils.PrintErrorMsg("Missing output file.")
				os.Exit(1)
			}
		}

		numWorkers, _ := cmd.Flags().GetInt("workers")
		getCommitsFlag, _ := cmd.Flags().GetBool("commits")

		// Create a channel to send work to the workers with a buffer size of length gitRepos
		jobs := make(chan string, len(gitRepos))
		repositories := make(chan utils.RepositoryJob, len(gitRepos))

		// Create a channel to indicate when the go routines have finished
		done := make(chan bool)

		var config utils.Config
		// Initialize the repositories map
		config.Repositories = make(map[string]utils.Repository)
		// Iterate over the numWorkers
		for i := 0; i < numWorkers; i++ {
			go func() {
				for repoPath := range jobs {
					var repoPathName string
					if repoPath == "." {
						absPath, _ := filepath.Abs(repoPath)
						repoPathName = filepath.Base(absPath)
					} else {
						repoPathName = filepath.Base(repoPath)
					}
					repo := utils.ParseRepositoryInfo(repoPath, getCommitsFlag)
					repositories <- utils.RepositoryJob{RepoPath: repoPathName, Repo: repo}
				}
				done <- true
			}()
		}
		// Send each git repository path to the jobs channel
		for _, repoPath := range gitRepos {
			jobs <- repoPath
		}
		close(jobs) // Close channel to signal no more work will be sent

		// wait for all goroutines to finish
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(repositories)

		for repoResult := range repositories {
			config.Repositories[repoResult.RepoPath] = repoResult.Repo
		}
		yamlData, _ := yaml.Marshal(&config)
		if visualizeOutput {
			fmt.Println(string(yamlData))
		}
		if !skipOutputFile {
			err := os.WriteFile(filePath, yamlData, 0644)
			if err != nil {
				utils.PrintErrorMsg("Failed to export repositories to yaml file.")
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
	exportCmd.Flags().StringP("output", "o", "", "Path to output `.repos` file")
	exportCmd.Flags().BoolP("commits", "c", false, "Export repositories hashes instead of branches")
	exportCmd.Flags().BoolP("visualize", "v", false, "Show the information to be stored in the output file")
}
