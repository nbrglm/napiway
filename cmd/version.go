/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/nbrglm/napiway/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the NApiWay generator.",
	Long:  `Prints the version of the NApiWay generator. This can be used to verify which version of the generator you are using, especially when troubleshooting or reporting issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Lookup("short").Value.String() == "true" {
			fmt.Println(version.Version)
		} else {
			fmt.Printf("NApiWay Generator version: %s\n", version.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionCmd.Flags().BoolP("short", "s", false, "Print only the version number without additional text")
}
