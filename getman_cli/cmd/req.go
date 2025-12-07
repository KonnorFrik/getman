/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reqCmd represents the req command
var reqCmd = &cobra.Command{
	Use:   "req",
	Short: "Work with requests",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("req called")
	},
}

func init() {
	rootCmd.AddCommand(reqCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reqCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reqCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
