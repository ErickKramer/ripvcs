package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/yaml"
)

type RepositoryJob struct {
	RepoPath string
	Repo     Repository
}

type Repository struct {
	Type    string `yaml:"type"`
	URL     string `yaml:"url"`
	Version string `yaml:"version,omitempty"`
}

type Config struct {
	Repositories map[string]Repository `yaml:"repositories"`
}

// IsReposFileValid Check if given filePath exists and if has .repos suffix
func IsReposFileValid(filePath string) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("Error: File does not exist!")
	}

	if !strings.HasSuffix(filePath, ".repos") {
		return errors.New("Error: File given does not have a valid .repos extension!")
	}
	return nil
}

// ParseReposFile Load data from a given .repos file
func ParseReposFile(filePath string) (*Config, error) {
	errValid := IsReposFileValid(filePath)
	if errValid != nil {
		return nil, errValid
	}
	yamlFile, err := os.ReadFile(filePath)

	if err != nil {
		errorMsg := "Failed to read .repos file"
		// fmt.Printf("%s: %s\n", errorMsg, err)
		return nil, errors.New(errorMsg)
	}

	// parse YAML content
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		errorMsg := "Failed to parse .repos file"
		// fmt.Printf("%s: %s\n", errorMsg, err)
		return nil, errors.New(errorMsg)
	}

	if len(config.Repositories) == 0 {
		errorMsg := "Empty .repos file given"
		// fmt.Printf("%s: %s\n", errorMsg, err)
		return nil, errors.New(errorMsg)
	}
	return &config, nil

}

// FindReposFiles Search .repos files in a given path
func FindReposFiles(rootPath string) ([]string, error) {
	var foundReposFiles []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".repos" {
			if err != nil {
				return err
			}
			foundReposFiles = append(foundReposFiles, path)
		}
		return nil
	})

	return foundReposFiles, err
}

// FindDirectory Search for a targetDir given a rootPath
func FindDirectory(rootPath string, targetDir string) (string, error) {
	if len(rootPath) == 0 {
		return "", errors.New("Empty rootPath given")
	}
	if len(targetDir) == 0 {
		return "", errors.New("Empty targetDir given")
	}
	if rootInfo, err := os.Stat(rootPath); err != nil || !rootInfo.IsDir() {
		return "", err
	}
	if _, err := os.Stat(targetDir); err == nil {
		return "", errors.New("targetDir is a Path!. Expected just a name.")
	}

	var dirPath string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == targetDir {
			dirPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

// ParseRepositoryInfo Create a Repository object containing the given repository info
func ParseRepositoryInfo(repoPath string, useCommit bool) Repository {
	var repository Repository
	if !IsGitRepository(repoPath) {
		return repository
	}
	repository.Type = "git"
	repository.URL = GetGitRemoteURL(repoPath)
	if useCommit {
		repository.Version = GetGitCommitSha(repoPath)
	} else {
		repository.Version = GetGitBranch(repoPath)
	}
	return repository
}

func GetRepoPath(repoName string) string {
	repoNameInfo, err := os.Stat(repoName)

	if err == nil {
		if !repoNameInfo.IsDir() {
			PrintErrorMsg(fmt.Sprintf("%s is not a directory\n", repoName))
			os.Exit(1)
		}
		return repoName
	}

	if !os.IsNotExist(err) {
		PrintErrorMsg(fmt.Sprintf("Error checking repository: %s\n", repoName))
		os.Exit(1)
	}

	foundRepoPath, findErr := FindDirectory(".", repoName)
	if findErr != nil {
		PrintErrorMsg(fmt.Sprintf("Failed to find directory named %s. Error: %s\n", repoName, findErr))
		os.Exit(1)
	}
	return foundRepoPath
}
