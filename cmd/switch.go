/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch <repo name | path>",
	Short: "Switch repository version",
	Long: `Switch repository version.

It allows to easily run Git switch operation on the given repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoName string
		if len(args) == 0 {
			fmt.Println("Repository Name or Path not given")
			os.Exit(1)
		}
		repoName = args[0]

		var repoPath string
		// check if input is not a path
		if repoNameInfo, err := os.Stat(repoName); err != nil {
			if os.IsNotExist(err) {
				foundRepoPath, findErr := utils.FindDirectory(".", repoName)
				if findErr != nil {
					fmt.Printf("Failed to find directory named %s. Error: %s\n", repoPath, findErr)
					os.Exit(1)
				}
				repoPath = foundRepoPath
			} else {
				fmt.Printf("Error checking repoName: %s\n", err)
				os.Exit(1)
			}
		} else if !repoNameInfo.IsDir() {
			fmt.Printf("%s is not a directory\n", repoName)
			os.Exit(1)
		} else {
			repoPath = repoName
		}

		if !utils.IsGitRepository(repoPath) {
			fmt.Println("Directory given is not a git repository")
			os.Exit(1)
		}
		createBranch, _ := cmd.Flags().GetBool("create")
		detachHead, _ := cmd.Flags().GetBool("detach")
		branch, _ := cmd.Flags().GetString("branch")
		utils.PrintGitSwitch(repoPath, branch, createBranch, detachHead)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().BoolP("create", "c", false, "Create and switch to a new branch")
	switchCmd.Flags().BoolP("detach", "d", false, "Detach HEAD at named commit or tag")
	switchCmd.Flags().StringP("branch", "b", "", "Version (branch, commit, or tag) to switch to")
}
