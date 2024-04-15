/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("import called")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("input", "i", "")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	importCmd.Flags().StringP("input", "i", "", "Path to input `.repos` file")
	importCmd.Flags().BoolP("recursive", "r", false, "Recursively search of other `.repos` file in the cloned repositories")
	importCmd.Flags().BoolP("skip-if-existing", "s", false, "Skip existing repositories")
	importCmd.Flags().StringP("path", "p", "", "Path where to clone the repositories.")
}
