package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/jesseduffield/yaml"
)

// Define a struct to hold the key and the repository
type RepositoryJob struct {
	DirName string
	Repo    Repository
}

type Repository struct {
	Type    string `yaml:"type"`
	URL     string `yaml:"url"`
	Version string `yaml:"version,omitempty"`
}

type Config struct {
	Repositories map[string]Repository `yaml:"repositories"`
}

func IsReposFileValid(filePath string) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("Error: File does not exist!")
	}

	if !strings.HasSuffix(filePath, ".repos") {
		return errors.New("Error: File given does not have a valid .repos extension!")
	}
	return nil
}

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
