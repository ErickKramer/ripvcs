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
)

var importCmd = &cobra.Command{
	Use:   "import <optional path>",
	Short: "Import repositories listed in the given .repos file",
	Long: `Import repositories listed in the given .repos file

The repositories are cloned in the given path or in the current path.

It supports recursively searching for any other .repos file found at each
import cycle.`,
	Run: func(cmd *cobra.Command, args []string) {
		var cloningPath string
		if len(args) == 0 {
			cloningPath = "."
		} else {
			cloningPath = args[0]
		}

		// Get arguments
		filePath, _ := cmd.Flags().GetString("input")
		recursiveFlag, _ := cmd.Flags().GetBool("recusive")
		skipExisting, _ := cmd.Flags().GetBool("skip-if-existing")
		depthRecursive, _ := cmd.Flags().GetInt("depth-recursive")
		numWorkers, _ := cmd.Flags().GetInt("workers")

		// Import repository files in the given file
		validFile := singleCloneSweep(cloningPath, filePath, numWorkers, skipExisting)
		if !recursiveFlag {
			if !validFile {
				os.Exit(1)
			}
			os.Exit(0)
		}
		nestedImportClones(cloningPath, filePath, depthRecursive, numWorkers, skipExisting)

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("input", "i", "", "Path to input `.repos` file")
	importCmd.Flags().BoolP("recursive", "r", false, "Recursively search of other `.repos` file in the cloned repositories")
	importCmd.Flags().IntP("depth-recursive", "d", -1, "Regulates how many levels the recursive dependencies would be cloned.")
	importCmd.Flags().BoolP("skip-if-existing", "s", false, "Skip existing repositories")
	importCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
}

func singleCloneSweep(root string, filePath string, numWorkers int, skipExisting bool) bool {
	utils.PrintSection(fmt.Sprintf("Importing from %s", filePath))
	utils.PrintSeparator()
	config, err := utils.ParseReposFile(filePath)
	if err != nil {
		fmt.Printf("Invalid file given {%s}. %s\n", filePath, err)
		return false
	}
	// Create a channel to send work to the workers with a buffer size of length gitRepos
	jobs := make(chan utils.RepositoryJob, len(config.Repositories))
	// Create channel to collect results
	results := make(chan bool, len(config.Repositories))
	// Create a channel to indicate when the go routines have finished
	done := make(chan bool)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for job := range jobs {
				if job.Repo.Type != "git" {
					utils.PrintRepoEntry(job.DirName, "")
					utils.PrintErrorMsg("Unsupported repository type.\n")
					results <- false
				} else {
					success := utils.PrintGitClone(job.Repo.URL, job.Repo.Version, job.DirName, skipExisting, false)
					results <- success
				}
			}
			done <- true
		}()
	}

	for dirName, repo := range config.Repositories {
		jobs <- utils.RepositoryJob{DirName: filepath.Join(root, dirName), Repo: repo}
	}
	close(jobs)
	// wait for all goroutines to finish
	for i := 0; i < numWorkers; i++ {
		<-done
	}
	close(results)

	validFile := true
	for result := range results {
		if !result {
			validFile = false
			break
		}
	}
	return validFile
}

func nestedImportClones(cloningPath string, initialFilePath string, depthRecursive int, numWorkers int, skipExisting bool) {
	// Recursively import .repos files found
	clonedReposFiles := map[string]bool{initialFilePath: true}
	validFiles := true
	cloneSweepCounter := 0

	for {
		// Check if recursion level has been reached
		if depthRecursive != -1 && cloneSweepCounter >= depthRecursive {
			break
		}

		// Find .repos file to clone
		foundReposFiles, err := utils.FindReposFiles(cloningPath)
		if err != nil || len(foundReposFiles) == 0 {
			break
		}

		// Get dependencies to clone
		newReposFileFound := false
		for _, filePathToClone := range foundReposFiles {
			if _, ok := clonedReposFiles[filePathToClone]; !ok {
				validFiles = singleCloneSweep(cloningPath, filePathToClone, numWorkers, skipExisting)
				clonedReposFiles[filePathToClone] = true
				newReposFileFound = true
				if !validFiles {
					fmt.Println("Encountered errors while importing file")
				}
			}
		}
		if !newReposFileFound {
			break
		}
		cloneSweepCounter++
	}

	if !validFiles {
		os.Exit(1)
	}
}
