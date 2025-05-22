/*
Copyright © 2024 Erick Kramer <erickkramer@gmail.com>
*/
package cmd

import (
	"fmt"
	"ripvcs/utils"

	"github.com/spf13/cobra"
)

var (
	// These variables should be set using ldflags during the build
	Version   = "v1.0.2"
	Commit    = "071f3a7"
	BuildDate = "22.05.2025"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSection("ripvcs (rv)")
		utils.PrintSeparator()
		fmt.Printf("%sVersion:%s %s\n", utils.BlueColor, utils.ResetColor, Version)
		fmt.Printf("%sCommit:%s %s\n", utils.BlueColor, utils.ResetColor, Commit)
		fmt.Printf("%sBuild Date:%s %s\n", utils.BlueColor, utils.ResetColor, BuildDate)
		utils.PrintSeparator()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
