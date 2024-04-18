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

func RunGitCmd(path string, gitCmd string, envConfig []string, args ...string) (string, error) {
	cmdArgs := append([]string{"-c", "color.ui=always", gitCmd}, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = append(os.Environ(), envConfig...)
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetGitStatus Execute git status in a given path
func GetGitStatus(path string) string {
	output, err := RunGitCmd(path, "status", nil)
	if err != nil {
		fmt.Printf("Failed to check Git status of %s. Error: %s", path, err)
	}
	return output
}

func IsGitURLValid(url string, branch string, enablePrompt bool) bool {
	var envConfig []string
	if enablePrompt {
		envConfig = []string{"GIT_TERMINAL_PROMPT=1"}
	} else {
		envConfig = []string{"GIT_TERMINAL_PROMPT=0"}
	}

	urlArgs := []string{url, branch}
	output, err := RunGitCmd(".", "ls-remote", envConfig, urlArgs...)
	if err != nil || len(output) == 0 {
		// fmt.Printf("Failed to check Git URL %s. Error: %s", url, err)
		return false
	}
	return true
}

func GetGitLog(path string, oneline bool, numCommits int) string {
	var cmdArgs []string

	if oneline {
		cmdArgs = []string{"-n", strconv.Itoa(numCommits), "--oneline"}
	} else {
		cmdArgs = []string{"-n", strconv.Itoa(numCommits)}
	}

	output, err := RunGitCmd(path, "log", nil, cmdArgs...)
	if err != nil {
		fmt.Printf("Failed to check Git log of %s. Error: %s", path, err)
	}
	return output
}

func PrintHelper(path string) {

	blueColor := "\033[38;2;137;180;250m"
	resetColor := "\033[0m"

	fmt.Printf("%s=== %s ===%s\n", blueColor, path, resetColor)
}

func PrintGitLog(path string, oneline bool, numCommits int) {
	repoLogs := GetGitLog(path, oneline, numCommits)

	PrintHelper(path)
	fmt.Print(string(repoLogs))
}

// PrintGitStatus Send the git status of a path to stdout with color codes.
func PrintGitStatus(path string, skipEmpty bool) {
	repoStatus := GetGitStatus(path)

	if skipEmpty && strings.Contains(repoStatus, "working tree clean") {
		return
	}

	PrintHelper(path)
	fmt.Print(string(repoStatus))
}

func PrintCheckGit(path string, url string, version string, enablePrompt bool) bool {
	redColor := "\033[38;2;255;0;0m"
	resetColor := "\033[0m"
	var checkMsg string
	var isURLValid bool
	if !IsGitURLValid(url, version, enablePrompt) {
		checkMsg = fmt.Sprintf("%sFailed to contact git repository '%s' with version '%s%s'\n", redColor, url, version, resetColor)
		isURLValid = false
	} else {
		checkMsg = fmt.Sprintf("Successfully contact git repository '%s' with version '%s'\n", url, version)
		isURLValid = true
	}
	PrintHelper(path)
	fmt.Printf(checkMsg)
	return isURLValid
}
