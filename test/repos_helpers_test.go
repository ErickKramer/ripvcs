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
	_, err := utils.ParseReposFile("./valid_example.repos")
	if err != nil {
		t.Errorf("Expected to report a valid file")
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
	foundReposFiles, err = utils.FindReposFiles("/tmp")

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

	repoPath, err = utils.FindDirectory("", "sadsd")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on empty rootPath")
	}
	repoPath, err = utils.FindDirectory("/tmp", "")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on empty targetDir")
	}
	repoPath, err = utils.FindDirectory("/sdasd", "")
	if err == nil {
		t.Errorf("Expected to failed to find directory, based on nonexisting rootPath")
	}
	repoPath, err = utils.FindDirectory("/tmp", "/tmp/testdata/")
	if err == nil {
		t.Errorf("Expected to failed to find directory, targetDir being a path")
	}
	err = os.RemoveAll("/tmp/testdata")
	if err != nil {
		panic(err)
	}

}
