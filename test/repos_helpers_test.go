package test

import (
	"os"
	"ripvcs/utils"
	"testing"
)

func TestNonReposFile(t *testing.T) {
	err := os.WriteFile("/tmp/wrong_extension.txt", []byte{}, 0644)
	if err != nil {
		panic(err)
	}
	if utils.IsReposFileValid("/tmp/wrong_extension.txt") {
		t.Errorf("Expected to report file with wrong extension")
	}
	if utils.IsReposFileValid("missing_file.repos") {
		t.Errorf("Expected to report non-existing file")
	}
}

func TestReposFile(t *testing.T) {
	if !utils.IsReposFileValid("./valid_example.repos") {
		t.Errorf("Expected to report a valid file")
	}
}
