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
		// var repoName string
		if len(args) == 0 {
			utils.PrintErrorMsg("Repository Name or Path not given\n")
			os.Exit(1)
		}
		repoPath := utils.GetRepoPath(args[0])

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
