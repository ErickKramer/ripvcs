package test

import (
	"fmt"
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

func TestCheckGitUrl(t *testing.T) {
	if utils.IsGitURLValid("https://github.com/ros2/demosasdasd.git", "rolling", false) {
		t.Errorf("Expected to return invalid URL")
	}
	if !utils.IsGitURLValid("https://github.com/ros2/demos.git", "rolling", false) {
		t.Errorf("Expected to return invalid URL")
	}
}

func TestCloneGitRepo(t *testing.T) {
	if utils.GitClone("https://github.com/ros2/demos.git", "rolling", "/tmp/testdata/demos", false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/ros2cli", "", "/tmp/testdata/ros2cli", false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/sadasdasd.git", "", "/tmp/testdata/sdasda", false, false) != utils.FailedClone {
		t.Errorf("Expected to fail to clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/demos.git", "", "/tmp/testdata/demos", true, false) != utils.SkippedClone {
		t.Errorf("Expected to skip to clone git repository")
	}
}

func TestGitSwitch(t *testing.T) {
	repoPath := "/tmp/testdata/switch_test"
	if utils.GitClone("https://github.com/ros2/demos.git", "rolling", repoPath, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	output, err := utils.GitSwitch(repoPath, "humble", false, false)
	if err != nil {
		t.Errorf("Expected to successfully to switch to a branch. Error %s", err)
	}

	output, err = utils.GitSwitch(repoPath, "nonexisting", false, false)
	if err == nil {
		t.Errorf("Expected to fail to switch to a nonexisting branch.\nError %s", err)
	}
	fmt.Println(output)
	output, err = utils.GitSwitch(repoPath, "nonexisting", true, false)
	if err != nil {
		t.Errorf("Expected to successfully to create a new branch.\nError %s", err)
	}
	output, err = utils.GitSwitch(repoPath, "0.34.0", false, true)
	if err != nil {
		t.Errorf("Expected to successfully to switch to a tag.\nError %s", err)
	}
}
