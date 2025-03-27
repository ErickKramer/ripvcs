package test

import (
	"os"
	"os/exec"
	"ripvcs/utils"
	"strings"
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
	if err != nil {
		panic(err)
	}
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
	if utils.GetGitStatus("/tmp/testdata/valid_repo", false) == "" {
		t.Errorf("Failed to check status of a valid repository")
	}
	if utils.GetGitStatus("/tmp/testdata/valid_repo", true) == "" {
		t.Errorf("Failed to check status of a valid repository")
	}
}

func TestGetGitBranch(t *testing.T) {
	testingBranch := "jazzy"
	repoPath := "/tmp/testdata/demos_branch"
	if utils.GitClone("https://github.com/ros2/demos.git", testingBranch, repoPath, true, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GetGitBranch(repoPath) != testingBranch {
		t.Errorf("Failed to get main branch for valid git repository")
	}
	testingTag := "0.34.0"
	if utils.GitClone("https://github.com/ros2/demos.git", testingTag, repoPath, true, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	obtainedTag := utils.GetGitBranch(repoPath)
	if obtainedTag != testingTag {
		t.Errorf("Failed to get tag for the cloned repository. Got %s", obtainedTag)
	}
}

func TestGetGitCommitSha(t *testing.T) {
	testingSha := "839b622bc40ec62307d6ba0615adb9b8bd1cbc30"
	repoPath := "/tmp/testdata/demos_sha"
	if utils.GitClone("https://github.com/ros2/demos.git", testingSha, repoPath, false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GetGitCommitSha(repoPath) != testingSha {
		t.Errorf("Failed to get commit sha of the cloned git repository")
	}
}

func TestGetGitRemoteURL(t *testing.T) {
	repoPath := "/tmp/testdata/demos_url"
	remoteUrl := "https://github.com/ros2/demos.git"
	if utils.GitClone(remoteUrl, "", repoPath, false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GetGitRemoteURL(repoPath) != remoteUrl {
		t.Errorf("Failed to get remote URL for the git repository")
	}
}

func TestGitPull(t *testing.T) {
	repoPath := "/tmp/testdata/demos_pull"
	if utils.GitClone("https://github.com/ros2/demos.git", "rolling", repoPath, false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	msg := utils.PullGitRepo(repoPath)
	if strings.TrimSpace(msg) != "Already up to date." {
		t.Errorf("Failed to pull valid repository. Obtained %s", msg)
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
	if valid, err := utils.IsGitURLValid("https://github.com/ros2/demosasdasd.git", "rolling", false); valid {
		t.Errorf("Expected to return invalid URL. Error %v", err)
	}
	if valid, err := utils.IsGitURLValid("https://github.com/ros2/demos.git", "rolling", false); !valid {
		t.Errorf("Expected to return valid URL given a branch. Error %v", err)
	}
	if valid, err := utils.IsGitURLValid("https://github.com/ros2/demos.git", "", false); !valid {
		t.Errorf("Expected to return valid URL given no branch. Error %v", err)
	}
	if valid, err := utils.IsGitURLValid("https://github.com/ros2/demos.git", "0.34.0", false); !valid {
		t.Errorf("Expected to return valid URL given a tag. Error %v", err)
	}
	if valid, err := utils.IsGitURLValid("https://github.com/ros2/demos.git", "839b622bc40ec62307d6ba0615adb9b8bd1cbc30", false); valid {
		t.Errorf("Expected to return invalid URL given a commit SHA. Error %v", err)
	}
}

func TestIsValidCommitSha(t *testing.T) {
	if !utils.IsValidSha("e69de29bb2d1d6434b8b29ae775ad8c2e48c5391") {
		t.Errorf("Expected to return valid SHA")
	}
	if !utils.IsValidSha("839b622") {
		t.Errorf("Expected to return valid SHA")
	}
	if utils.IsValidSha("INVALIDSHA123456789012345678901234567890") {
		t.Errorf("Expected to return invalid SHA")
	}
	if utils.IsValidSha("11111111111111111111111111111111111111111") {
		t.Errorf("Expected to return invalid SHA. Invalid length (41 chars)")
	}
	if utils.IsValidSha("") {
		t.Errorf("Expected to return invalid SHA. Invalid length (0 chars)")
	}
	if utils.IsValidSha("999999") {
		t.Errorf("Expected to return invalid SHA. Invalid length (6 chars)")
	}
}

func TestCloneGitRepo(t *testing.T) {
	repoPath := "/tmp/testdata/demos_clone"
	if utils.GitClone("https://github.com/ros2/demos.git", "rolling", repoPath, false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/ros2cli", "", "/tmp/testdata/ros2cli", false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/sadasdasd.git", "", "/tmp/testdata/sdasda", false, false, false, false) != utils.FailedClone {
		t.Errorf("Expected to fail to clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/demos.git", "", repoPath, false, false, false, false) != utils.SkippedClone {
		t.Errorf("Expected to skip to clone git repository")
	}
	if utils.GitClone("https://github.com/ros2/demos.git", "", repoPath, true, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to overwrite found git repository")
	}
	if utils.GitClone("https://github.com/ros2/demos.git", "", repoPath, true, true, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully to clone git repository with shallow enabled")
	}
	count, err := utils.RunGitCmd(repoPath, "rev-list", nil, []string{"--all", "--count"}...)
	if err != nil || strings.TrimSpace(count) != "1" {
		t.Errorf("Expected to have a shallow clone of the git repository")
	}

	testingSha := "839b622bc40ec62307d6ba0615adb9b8bd1cbc30"
	if utils.GitClone("https://github.com/ros2/demos.git", testingSha, repoPath, true, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository given a SHA")
	}
}

func TestGitSwitch(t *testing.T) {
	repoPath := "/tmp/testdata/switch_test"
	if utils.GitClone("https://github.com/ros2/demos.git", "rolling", repoPath, false, false, false, false) != utils.SuccessfullClone {
		t.Errorf("Expected to successfully clone git repository")
	}
	_, err := utils.GitSwitch(repoPath, "humble", false, false)
	if err != nil {
		t.Errorf("Expected to successfully to switch to a branch. Error %s", err)
	}

	_, err = utils.GitSwitch(repoPath, "nonexisting", false, false)
	if err == nil {
		t.Errorf("Expected to fail to switch to a nonexisting branch.\nError %s", err)
	}
	_, err = utils.GitSwitch(repoPath, "nonexisting", true, false)
	if err != nil {
		t.Errorf("Expected to successfully to create a new branch.\nError %s", err)
	}
	_, err = utils.GitSwitch(repoPath, "0.34.0", false, true)
	if err != nil {
		t.Errorf("Expected to successfully to switch to a tag.\nError %s", err)
	}
}
