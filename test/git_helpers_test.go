package test

import (
	"os"
	"ripvcs/utils"
	"testing"
)

func TestMain(m *testing.M) {
	// Create temporary directories and files for testing
	createTestFiles()

	// Run tests
	exitVal := m.Run()

	cleanupTestFiles()

	// Exit with the appropriate exit code
	os.Exit(exitVal)
}

func createTestFiles() {
	// Create root directory
	err := os.MkdirAll("testdata/valid_repo/.git", 0755)
	if err != nil {
		panic(err)
	}

	// Create nested directories and .git repository
	err = os.MkdirAll("testdata/normal_dir/another_repo/.git", 0755)
	if err != nil {
		panic(err)
	}
}

func cleanupTestFiles() {
	err := os.RemoveAll("testdata")
	if err != nil {
		panic(err)
	}
}

func TestIsGitRepository(t *testing.T) {
	if !utils.IsGitRepository("./testdata/valid_repo") {
		t.Errorf("Expected ./valid_repo to be a Git repository")
	}

	if utils.IsGitRepository("./testdata/normal_dir") {
		t.Errorf("Expected ./normal_dir to not be a Git repository")
	}
}

func TestFindGitRepos(t *testing.T) {
	foundRepos := utils.FindGitRepositories("testdata")
	if len(foundRepos) != 2 {
		t.Errorf("Expected two git repositories relative to this test file, but got %v", len(foundRepos))
	}
}
