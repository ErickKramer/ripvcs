package test

import (
	"os"
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
}

func TestFindingReposFiles(t *testing.T) {
	foundReposFiles, err := utils.FindReposFiles(".")

	if err != nil || len(foundReposFiles) == 0 {
		t.Errorf("Expected to find at least one .repos file %v", err)
	}
	foundReposFiles, err = utils.FindReposFiles("/tmp/stuff")

	if err != nil || len(foundReposFiles) != 0 {
		t.Errorf("Expected to not find any .repos file %v", err)
	}
}
