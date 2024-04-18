// utils/git_helpers.go

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// IsGitRepository checks if a directory is a git repository
func IsGitRepository(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// FindGitRepositories Get a slice of all the found git repositories at the given root
func FindGitRepositories(root string) []string {
	var gitRepos []string

	// Use an anonymous function to  check each file found relative to the given root
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return any error encountered during walking
		}
		if info.IsDir() && IsGitRepository(path) {
			gitRepos = append(gitRepos, path)
		}
		return nil // Continue walking
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return gitRepos
}

func RunGitCmd(path string, gitCmd string, args ...string) string {
	cmdArgs := append([]string{"-c", "color.ui=always", gitCmd}, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running 'git %s %s' at %s: %s\n", gitCmd, strings.Join(args, " "), path, err)
		return ""
	}
	return string(output)
}

// GetGitStatus Execute git status in a given path
func GetGitStatus(path string) string {
	return RunGitCmd(path, "status")
}

func GetGitLog(path string, oneline bool, numCommits int) string {
	var cmdArgs []string

	if oneline {
		cmdArgs = []string{"-n", strconv.Itoa(numCommits), "--oneline"}
	} else {
		cmdArgs = []string{"-n", strconv.Itoa(numCommits)}
	}

	repoLogs := RunGitCmd(path, "log", cmdArgs...)
	return repoLogs
}

func PrintGitLog(path string, oneline bool, numCommits int) {
	repoLogs := GetGitLog(path, oneline, numCommits)

	blueColor := "\033[38;2;137;180;250m"
	resetColor := "\033[0m"

	fmt.Printf("%s=== %s ===%s\n", blueColor, path, resetColor)
	fmt.Print(string(repoLogs))
}

// PrintGitStatus Send the git status of a path to stdout with color codes.
func PrintGitStatus(path string, skipEmpty bool) {
	repoStatus := GetGitStatus(path)

	blueColor := "\033[38;2;137;180;250m"
	resetColor := "\033[0m"

	if skipEmpty && strings.Contains(repoStatus, "working tree clean") {
		return
	}

	fmt.Printf("%s=== %s ===%s\n", blueColor, path, resetColor)
	fmt.Print(string(repoStatus))
}
