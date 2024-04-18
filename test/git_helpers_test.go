package test

import (
	"os"
	"os/exec"
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
	path := "/tmp/testdata/valid_repo/"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// Create nested directories and .git repository
	path = "/tmp/testdata/normal_dir/another_repo/"
	err = os.MkdirAll(path, 0755)
	cmd = exec.Command("git", "init")
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
}

func cleanupTestFiles() {
	err := os.RemoveAll("/tmp/testdata")
	if err != nil {
		panic(err)
	}
}

func TestIsGitRepository(t *testing.T) {
	if !utils.IsGitRepository("/tmp/testdata/valid_repo") {
		t.Errorf("Expected ./valid_repo to be a Git repository")
	}

	if utils.IsGitRepository("/tmp/testdata/normal_dir") {
		t.Errorf("Expected ./normal_dir to not be a Git repository")
	}
}

func TestFindGitRepos(t *testing.T) {
	foundRepos := utils.FindGitRepositories("/tmp/testdata")
	if len(foundRepos) != 2 {
		t.Errorf("Expected two git repositories relative to this test file, but got %v", len(foundRepos))
	}
}

func TestGitStatus(t *testing.T) {
	if utils.GetGitStatus("/tmp/testdata/valid_repo") == "" {
		t.Errorf("Failed to check status of a valid repository")
	}
}

func TestGitLog(t *testing.T) {
	if utils.GetGitLog(".", false, 5) == "" {
		t.Errorf("Failed to check logs of a valid repository")
	}

	if utils.GetGitLog(".", true, 5) == "" {
		t.Errorf("Failed to check logs of a valid repository")
	}
}
