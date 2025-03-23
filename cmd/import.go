/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"ripvcs/utils"
	"strings"

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
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		numRetries, _ := cmd.Flags().GetInt("retry")
		overwriteExisting, _ := cmd.Flags().GetBool("force")
		shallowClone, _ := cmd.Flags().GetBool("shallowClone")
		depthRecursive, _ := cmd.Flags().GetInt("depth-recursive")
		numWorkers, _ := cmd.Flags().GetInt("workers")
		excludeList, _ := cmd.Flags().GetStringSlice("exclude")

		// Import repository files in the given file
		validFile := singleCloneSweep(cloningPath, filePath, numWorkers, overwriteExisting, shallowClone, numRetries)
		if !validFile {
			os.Exit(1)
		}
		if !recursiveFlag {
			os.Exit(0)
		}
		nestedImportClones(cloningPath, filePath, depthRecursive, numWorkers, overwriteExisting, shallowClone, numRetries, excludeList)

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().IntP("depth-recursive", "d", -1, "Regulates how many levels the recursive dependencies would be cloned.")
	importCmd.Flags().StringP("input", "i", "", "Path to input `.repos` file")
	importCmd.Flags().BoolP("recursive", "r", false, "Recursively search of other `.repos` file in the cloned repositories")
	importCmd.Flags().IntP("retry", "n", 2, "Number of attempts to import repositories")
	importCmd.Flags().BoolP("force", "f", false, "Force overwriting existing repositories")
	importCmd.Flags().BoolP("shallow", "l", false, "Clone repositories with a depth of 1")
	importCmd.Flags().IntP("workers", "w", 8, "Number of concurrent workers to use")
	importCmd.Flags().StringSliceP("exclude", "x", []string{}, "List of files or directories to exclude")
}

func singleCloneSweep(root string, filePath string, numWorkers int, overwriteExisting bool, shallowClone bool, numRetries int) bool {
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

	for range numWorkers {
		go func() {
			for job := range jobs {
				if job.Repo.Type != "git" {
					utils.PrintRepoEntry(job.RepoPath, "")
					utils.PrintErrorMsg(fmt.Sprintf("Unsupported repository type %s.\n", job.Repo.Type))
					results <- false
				} else {
					success := false
					for range numRetries {
						success = utils.PrintGitClone(job.Repo.URL, job.Repo.Version, job.RepoPath, overwriteExisting, shallowClone, false)
						if success {
							break
						}
					}
					results <- success
				}
			}
			done <- true
		}()
	}

	for dirName, repo := range config.Repositories {
		jobs <- utils.RepositoryJob{RepoPath: filepath.Join(root, dirName), Repo: repo}
	}
	close(jobs)
	// wait for all goroutines to finish
	for range numWorkers {
		<-done
	}
	close(results)

	validFile := true
	for result := range results {
		if !result {
			validFile = false
			fmt.Printf("Failed while cloning %s\n", filePath)
			break
		}
	}
	return validFile
}

func nestedImportClones(cloningPath string, initialFilePath string, depthRecursive int, numWorkers int, overwriteExisting bool, shallowClone bool, numRetries int, excludeList []string) {
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
			// Check if the file is in the exclude list
			exclude := false
			for _, excludePath := range excludeList {
				cleanExcludePath := filepath.Clean(excludePath)
				relPath, err := filepath.Rel(cloningPath, filePathToClone)
				if err != nil {
					continue
				}
				cleanRelPath := filepath.Clean(relPath)
				if cleanRelPath == cleanExcludePath || strings.HasPrefix(cleanRelPath, cleanExcludePath+string(os.PathSeparator)) {
					exclude = true
					break
				}
			}
			if exclude {
				utils.PrintRepoEntry(fmt.Sprintf("Excluding %s", filePathToClone), "")
				continue
			}

			if _, ok := clonedReposFiles[filePathToClone]; !ok {
				validFiles = singleCloneSweep(cloningPath, filePathToClone, numWorkers, overwriteExisting, shallowClone, numRetries)
				clonedReposFiles[filePathToClone] = true
				newReposFileFound = true
				if !validFiles {
					utils.PrintErrorMsg("Encountered errors while importing file")
					os.Exit(1)
				}
			}
		}
		if !newReposFileFound {
			break
		}
		cloneSweepCounter++
	}
}
