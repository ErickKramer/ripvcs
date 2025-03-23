package test

import (
	"os"
	"path/filepath"
	"ripvcs/utils"
	"testing"
)

func TestReposFile(t *testing.T) {
	err := os.WriteFile("/tmp/wrong_extension.txt", []byte{}, 0644)
	if err != nil {
		panic(err)
	}
	errValid := utils.IsReposFileValid("/tmp/wrong_extension.txt")
	if errValid == nil {
		t.Errorf("Expected to report file with wrong extension")
	}
	errValid = utils.IsReposFileValid("missing_file.repos")
	if errValid == nil {
		t.Errorf("Expected to report non-existing file")
	}
	errValid = utils.IsReposFileValid("./valid_example.repos")
	if errValid != nil {
		t.Errorf("Expected to report a valid file")
	}
}

func TestParsingReposFile(t *testing.T) {
	repo, err := utils.ParseReposFile("./valid_example.repos")
	if err != nil {
		t.Errorf("Expected to parse .repos file")
	}
	if len(repo.Repositories) != 8 {
		t.Errorf("The parsed repositories from .repos file do not match the expected values.")
	}
	expectedType := "git"
	if repoType := repo.Repositories["demos_rolling"].Type; repoType != expectedType {
		t.Errorf("Expected to have %s as repo type. Got: %s", expectedType, repoType)
	}
	expectedUrl := "https://github.com/ros2/demos.git"
	if repoUrl := repo.Repositories["demos_rolling"].URL; repoUrl != expectedUrl {
		t.Errorf("Expected to have %s as repo url. Got: %s", expectedUrl, repoUrl)
	}
	expectedVersion := "rolling"
	if repoVersion := repo.Repositories["demos_rolling"].Version; repoVersion != expectedVersion {
		t.Errorf("Expected to have %s as repo version. Got: %s", expectedVersion, repoVersion)
	}
	repo, err = utils.ParseReposFile("./valid_example.rosinstall")
	if err != nil {
		t.Errorf("Expected to parse .rosinstall file")
	}
	if len(repo.Repositories) != 2 {
		t.Errorf("The parsed repositories from .rosinstall do not match the expected values.")
	}
	expectedType = "git"
	if repoType := repo.Repositories["moveit_msgs"].Type; repoType != expectedType {
		t.Errorf("Expected to have %s as repo type. Got: %s", expectedType, repoType)
	}
	expectedUrl = "https://github.com/moveit/moveit_msgs.git"
	if repoUrl := repo.Repositories["moveit_msgs"].URL; repoUrl != expectedUrl {
		t.Errorf("Expected to have %s as repo url. Got: %s", expectedUrl, repoUrl)
	}
	expectedVersion = "master"
	if repoVersion := repo.Repositories["moveit_msgs"].Version; repoVersion != expectedVersion {
		t.Errorf("Expected to have %s as repo version. Got: %s", expectedVersion, repoVersion)
	}

	err = os.WriteFile("/tmp/empty_file.repos", []byte{}, 0644)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseReposFile("/tmp/empty_file.repos")
	if err == nil {
		t.Errorf("Expected to report empty file")
	}
	err = os.RemoveAll("/tmp/empty_file.repos")
	if err != nil {
		panic(err)
	}
}

func TestFindingReposFiles(t *testing.T) {
	foundReposFiles, err := utils.FindReposFiles(".")

	if err != nil || len(foundReposFiles) == 0 {
		t.Errorf("Expected to find at least one .repos file %v", err)
	}
	foundReposFiles, err = utils.FindReposFiles("../cmd/")

	if err != nil || len(foundReposFiles) != 0 {
		t.Errorf("Expected to not find any .repos file %v", err)
	}
}

func TestFindDirectory(t *testing.T) {
	// Create dummy dir
	path := "/tmp/testdata/valid_repo/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	repoPath, err := utils.FindDirectory("/tmp/testdata", "valid_repo")
	if err != nil {
		t.Errorf("Expected to find directory %v", err)
	}
	if filepath.Clean(repoPath) != filepath.Clean(path) {
		t.Errorf("Wrong directory found. Expected %v, found %v", path, repoPath)
	}

	_, err = utils.FindDirectory("", "sadsd")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on empty rootPath")
	}
	_, err = utils.FindDirectory("/tmp", "")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on empty targetDir")
	}
	_, err = utils.FindDirectory("/sdasd", "")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on nonexisting rootPath")
	}
	_, err = utils.FindDirectory("/tmp", "/tmp/testdata/")
	if err == nil {
		t.Errorf("Expected to failed to find directory, targetDir being a path")
	}
	err = os.RemoveAll("/tmp/testdata")
	if err != nil {
		panic(err)
	}
}

func TestParseRepositoryInfo(t *testing.T) {
	repository := utils.ParseRepositoryInfo("", false)
	if repository.Type != "" || repository.Version != "" || repository.URL != "" {
		t.Errorf("Expected to get an empty repository object")
	}
	repoPath := "/tmp/testdata/demos_parse"
	repoURL := "https://github.com/ros2/demos.git"
	repoVersion := "rolling"
	if utils.GitClone(repoURL, repoVersion, repoPath, true, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}

	repository = utils.ParseRepositoryInfo(repoPath, false)
	if repository.Type != "git" || repository.Version != repoVersion || repository.URL != repoURL {
		t.Errorf("Failed to properly parse the repository info using branch")
	}

	repoVersion = "839b622bc40ec62307d6ba0615adb9b8bd1cbc30"
	if utils.GitClone(repoURL, repoVersion, repoPath, true, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	repository = utils.ParseRepositoryInfo(repoPath, true)
	if repository.Type != "git" || repository.Version != repoVersion || repository.URL != repoURL {
		t.Errorf("Failed to properly parse the repository info using commit")
	}

	repoVersion = "0.34.0"
	if utils.GitClone(repoURL, repoVersion, repoPath, true, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	repository = utils.ParseRepositoryInfo(repoPath, false)
	if repository.Type != "git" || repository.Version != repoVersion || repository.URL != repoURL {
		t.Errorf("Failed to properly parse the repository info using tag")
	}
}
