// utils/git_helpers.go

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Create constant Error messages
const (
	SuccessfullClone = iota
	SkippedClone
	FailedClone
)

// IsGitRepository checks if a directory is a git repository
func IsGitRepository(dir string) bool {
	// FIXME: Check if dir is a directory
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

// RunGitCmd Helper method to execute a git command
func RunGitCmd(path string, gitCmd string, envConfig []string, args ...string) (string, error) {
	cmdArgs := append([]string{"-c", "color.ui=always", gitCmd}, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = append(os.Environ(), envConfig...)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetGitStatus Execute git status in a given path
func GetGitStatus(path string, plainStatus bool) string {
	var statusArgs []string
	if plainStatus {
		statusArgs = []string{"-sb"}
	}
	output, err := RunGitCmd(path, "status", nil, statusArgs...)
	if err != nil {
		fmt.Printf("Failed to check Git status of %s. Error: %s", path, err)
	}
	return output
}

// GetGitBranch Get current git branch in a given path
func GetGitBranch(path string) string {
	output, err := RunGitCmd(path, "branch", nil, "--show-current")
	if err != nil {
		fmt.Printf("Failed to get current Git branch of %s. Error: %s", path, err)
	}
	if output != "" {
		return strings.TrimSpace(output)
	}
	checkTagArgs := []string{"--points-at", "HEAD"}
	output, err = RunGitCmd(path, "tag", nil, checkTagArgs...)
	if err != nil {
		fmt.Printf("Failed to get current Git branch of %s. Error: %s", path, err)
	}
	return strings.TrimSpace(output)
}

func GetGitCommitSha(path string) string {
	cmdArgs := []string{"--verify", "HEAD"}
	output, err := RunGitCmd(path, "rev-parse", nil, cmdArgs...)
	if err != nil {
		fmt.Printf("Failed to get current Git commit of %s. Error: %s", path, err)
	}
	return strings.TrimSpace(output)
}

func GetGitRemoteURL(path string) string {
	cmdArgs := []string{"get-url", "origin"}
	output, err := RunGitCmd(path, "remote", nil, cmdArgs...)
	if err != nil {
		fmt.Printf("Failed to get URL for the origin remote of %s. Error: %s", path, err)
	}
	return strings.TrimSpace(output)
}

// PullGitRepo Execute git pull in a given path
func PullGitRepo(path string) string {
	output, err := RunGitCmd(path, "pull", nil)
	if err != nil {
		fmt.Printf("Failed to pull Git repository %s. Error: %s", path, err)
	}
	return output
}

// StashGitRepo Execute git stash in a given path
func StashGitRepo(path string, stashCmd string) string {
	output, err := RunGitCmd(path, "stash", nil, []string{stashCmd}...)
	if err != nil {
		fmt.Printf("Failed to run stash with %s Git repository %s. Error: %s", stashCmd, path, err)
	}
	return output
}

// SyncGitRepo Handle syncronization of a git repo
func SyncGitRepo(path string) string {
	output := StashGitRepo(path, "push")
	output += PullGitRepo(path)
	if StashGitRepo(path, "list") != "" {
		output += StashGitRepo(path, "pop")
	}
	return output
}

// IsGitURLValid Check if a git URL is reachable
func IsGitURLValid(url string, version string, enablePrompt bool) (bool, error) {
	var envConfig []string
	if enablePrompt {
		envConfig = []string{"GIT_TERMINAL_PROMPT=1"}
	} else {
		envConfig = []string{"GIT_TERMINAL_PROMPT=0"}
	}

	var urlArgs []string
	var output string
	var err error

	if IsValidSha(version) {
		err = fmt.Errorf("validation of URL given a commit SHA is currently not supported")
	} else {
		if version == "" {
			urlArgs = []string{url}
		} else {
			urlArgs = []string{url, version}
		}
		output, err = RunGitCmd(".", "ls-remote", envConfig, urlArgs...)
	}
	if err != nil || len(output) == 0 {
		return false, err
	}
	return true, nil
}

// GetGitLog Get logs for a given git repository
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

// GitSwitch Switch version for a given git repository
func GitSwitch(path string, branch string, createBranch bool, detachHead bool) (string, error) {

	cmdArgs := []string{}

	if detachHead {
		cmdArgs = append(cmdArgs, "--detach")
	} else if createBranch {
		cmdArgs = append(cmdArgs, "--create")
	}
	cmdArgs = append(cmdArgs, branch)

	output, err := RunGitCmd(path, "switch", nil, cmdArgs...)
	if err != nil {
		switchError := errors.New(fmt.Sprintf("Failed to switch branch of repository %s to %s. Error: %s", path, branch, err))
		return "", switchError
	}
	return output, nil
}

// IsValidSha Check if sha given is a valid SHA1
func IsValidSha(sha string) bool {
	shaRegex := regexp.MustCompile(`^[a-fA-F0-9]{7,40}$`)
	return shaRegex.MatchString(sha)

}

// GitClone Clone a given repository URL
func GitClone(url string, version string, clonePath string, overwriteExisting bool, shallowClone bool, enablePrompt bool) int {

	// Check if clonePath exists
	if _, err := os.Stat(clonePath); err == nil {
		if !overwriteExisting {
			return SkippedClone
		} else {
			// Remove existing clonePath
			if err := os.RemoveAll(clonePath); err != nil {
				fmt.Printf("Failed to remove existing cloning path %s. Error: %s\n", clonePath, err)
				panic(err)
			}
		}
	}

	var envConfig []string
	if enablePrompt {
		envConfig = []string{"GIT_TERMINAL_PROMPT=1"}
	} else {
		envConfig = []string{"GIT_TERMINAL_PROMPT=0"}
	}

	var cmdArgs []string

	versionIsSha := IsValidSha(version)
	if version == "" || versionIsSha {
		cmdArgs = []string{url, clonePath}
	} else {
		cmdArgs = []string{url, "--branch", version, clonePath}
	}

	if shallowClone {
		cmdArgs = append(cmdArgs, "--depth", "1")
	}
	if _, err := RunGitCmd(".", "clone", envConfig, cmdArgs...); err != nil {
		return FailedClone
	}

	if versionIsSha {
		if _, err := GitSwitch(clonePath, version, false, true); err != nil {
			return FailedClone
		}
	}

	return SuccessfullClone
}

// PrintGitLog Pretty print logs for a given git repository
func PrintGitLog(path string, oneline bool, numCommits int) {
	repoLogs := GetGitLog(path, oneline, numCommits)
	PrintRepoEntry(path, string(repoLogs))
}

// PrintGitStatus Pretty print status for a given git repository
func PrintGitStatus(path string, skipEmpty bool, plainStatus bool) {
	repoStatus := GetGitStatus(path, plainStatus)

	if plainStatus {
		if skipEmpty && strings.Count(repoStatus, "\n") <= 1 {
			return
		}
	} else {
		if skipEmpty && strings.Contains(repoStatus, "working tree clean") {
			return
		}
	}

	PrintRepoEntry(path, string(repoStatus))
}

// PrintGitPull Pretty print git pull output for a given git repository
func PrintGitPull(path string) {
	pullMsg := PullGitRepo(path)

	PrintRepoEntry(path, string(pullMsg))
}

// PrintGitSync Pretty print git sync output for a given git repository
func PrintGitSync(path string) {
	syncMsg := SyncGitRepo(path)

	PrintRepoEntry(path, string(syncMsg))
}

// PrintCheckGit Pretty print git url validation
func PrintCheckGit(path string, url string, version string, enablePrompt bool) bool {
	var checkMsg string
	var isURLValid bool
	if isURLValid, err := IsGitURLValid(url, version, enablePrompt); !isURLValid {
		checkMsg = fmt.Sprintf("%sFailed to contact git repository '%s' with version '%s'. Error: %v%s\n", RedColor, url, version, err, ResetColor)
	} else {
		checkMsg = fmt.Sprintf("Successfully contact git repository '%s' with version '%s'\n", url, version)
	}
	PrintRepoEntry(path, checkMsg)
	return isURLValid
}

// PrintGitClone Pretty print git clone
func PrintGitClone(url string, version string, path string, overwriteExisting bool, shallowClone bool, enablePrompt bool) bool {
	var cloneMsg string
	var cloneSuccessful bool
	statusClone := GitClone(url, version, path, overwriteExisting, shallowClone, enablePrompt)
	switch statusClone {
	case SuccessfullClone:
		cloneMsg = fmt.Sprintf("Successfully cloned git repository '%s' with version '%s'\n", url, version)
		cloneSuccessful = true
	case SkippedClone:
		cloneMsg = fmt.Sprintf("%sSkipped cloning existing git repository '%s'%s\n", OrangeColor, url, ResetColor)
		cloneSuccessful = true
	case FailedClone:
		cloneMsg = fmt.Sprintf("%sFailed to clone git repository '%s' with version '%s'%s\n", RedColor, url, version, ResetColor)
		cloneSuccessful = false
	default:
		panic("Unexpected behavior!")
	}
	PrintRepoEntry(path, cloneMsg)
	return cloneSuccessful
}

// PrintGitSwitch Pretty print git switch
func PrintGitSwitch(path string, branch string, createBranch bool, detachHead bool) bool {
	switchMsg, err := GitSwitch(path, branch, createBranch, detachHead)
	if err == nil {
		PrintRepoEntry(path, string(switchMsg))
		return true
	}
	errorMsg := fmt.Sprintf("%sError: '%s'%s\n", RedColor, err, ResetColor)
	PrintRepoEntry(path, string(errorMsg))
	return false
}
