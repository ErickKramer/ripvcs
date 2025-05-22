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
	"sync"

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
		recurseSubmodules, _ := cmd.Flags().GetBool("recurse-submodules")

		var hardCodedExcludeList = []string{}
		var clonedPaths []string

		// Import repository files in the given file
		validFile, hardCodedExcludeList, clonedPaths := singleCloneSweep(cloningPath, filePath, numWorkers, overwriteExisting, shallowClone, numRetries, recurseSubmodules)
		if !validFile {
			os.Exit(1)
		}
		if !recursiveFlag {
			os.Exit(0)
		}
		excludeList = append(excludeList, hardCodedExcludeList...)
		nestedImportClones(cloningPath, filePath, depthRecursive, numWorkers, overwriteExisting, shallowClone, numRetries, excludeList, recurseSubmodules, clonedPaths)

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
	importCmd.Flags().StringSliceP("exclude", "x", []string{}, "List of files and/or directories to exclude when performing a recursive import")
	importCmd.Flags().BoolP("recurse-submodules", "s", false, "Recursively clone submodules")
}

func singleCloneSweep(root string, filePath string, numWorkers int, overwriteExisting bool, shallowClone bool, numRetries int, recurseSubmodules bool) (bool, []string, []string) {
	utils.PrintSeparator()
	utils.PrintSection(fmt.Sprintf("Importing from %s", filePath))
	utils.PrintSeparator()
	config, err := utils.ParseReposFile(filePath)

	var allExcludes []string
	var clonedPaths []string

	if err != nil {
		utils.PrintErrorMsg(fmt.Sprintf("Invalid file given {%s}. %s\n", filePath, err))
		return false, allExcludes, clonedPaths
	}
	// Create a channel to send work to the workers with a buffer size of length gitRepos
	jobs := make(chan utils.RepositoryJob, len(config.Repositories))
	// Create channel to collect results
	results := make(chan bool, len(config.Repositories))
	// Create a channel to indicate when the go routines have finished
	done := make(chan bool)

	// Create mutex to handle excludeFilesChannel
	var excludeFilesMutex sync.Mutex

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
						success = utils.PrintGitClone(job.Repo.URL, job.Repo.Version, job.RepoPath, overwriteExisting, shallowClone, false, recurseSubmodules)
						if success {
							clonedPaths = append(clonedPaths, job.RepoPath)
							break
						}
					}
					results <- success
					// Expand excludeFilesChannel
					if len(job.Repo.Exclude) > 0 {
						excludeFilesMutex.Lock()
						allExcludes = append(allExcludes, job.Repo.Exclude...)
						excludeFilesMutex.Unlock()
					}
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
			utils.PrintErrorMsg(fmt.Sprintf("Failed while cloning %s\n", filePath))
			break
		}
	}

	return validFile, allExcludes, clonedPaths
}

func nestedImportClones(cloningPath string, initialFilePath string, depthRecursive int, numWorkers int, overwriteExisting bool, shallowClone bool, numRetries int, excludeList []string, recurseSubmodules bool, clonedPaths []string) {
	// Recursively import .repos files found
	clonedReposFiles := map[string]bool{initialFilePath: true}
	validFiles := true
	cloneSweepCounter := 0

	numPreviousFoundReposFiles := 0

	for {
		// Check if recursion level has been reached
		if depthRecursive != -1 && cloneSweepCounter >= depthRecursive {
			break
		}

		// Find .repos file to clone
		foundReposFiles, err := utils.FindReposFiles(cloningPath, clonedPaths)
		if err != nil || len(foundReposFiles) == 0 {
			break
		}

		if len(foundReposFiles) == numPreviousFoundReposFiles {
			break
		}
		numPreviousFoundReposFiles = len(foundReposFiles)

		// Get dependencies to clone
		newReposFileFound := false
		var hardCodedExcludeList = []string{}

		// FIXME: Find a simpler logic for this
		for _, filePathToClone := range foundReposFiles {
			// Check if the file is in the exclude list
			exclude := false

			// Initialize filePathToClone options
			filePathBase := filepath.Base(filePathToClone)
			filePathDir := filepath.Dir(filePathToClone)
			filePathParentDir := filepath.Base(filePathDir)

			for _, excludePath := range excludeList {
				excludeBase := filepath.Base(excludePath)

				// Check if exclude matches either:
				// 1. The full relative path
				// 2. The filename
				// 3. The parent directory
				if filePathBase == excludeBase || filePathParentDir == excludeBase || strings.HasPrefix(filePathToClone, excludePath) {
					exclude = true
					break
				}
			}

			if _, ok := clonedReposFiles[filePathToClone]; !ok {
				if exclude {
					utils.PrintSeparator()
					utils.PrintWarnMsg(fmt.Sprintf("Excluded cloning from '%s'\n", filePathToClone))
					clonedReposFiles[filePathToClone] = false
					continue
				}
				var newClonedPaths []string
				validFiles, hardCodedExcludeList, newClonedPaths = singleCloneSweep(cloningPath, filePathToClone, numWorkers, overwriteExisting, shallowClone, numRetries, recurseSubmodules)
				clonedReposFiles[filePathToClone] = true
				newReposFileFound = true
				clonedPaths = append(clonedPaths, newClonedPaths...)
				if !validFiles {
					utils.PrintErrorMsg("Encountered errors while importing file")
					os.Exit(1)
				}
				excludeList = append(excludeList, hardCodedExcludeList...)
			}
		}
		if !newReposFileFound {
			break
		}
		cloneSweepCounter++
	}
}
