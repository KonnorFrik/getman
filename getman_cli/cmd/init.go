/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a base directory with name 'dir'",
	Long: ``,
	Args: cobra.NoArgs,
	Run: _InitCmd,
}

func _InitCmd(cmd *cobra.Command, args []string) {
	if dirFlag == "" {
		PrintfCobraError(cmd, "Flag 'dir' cannot be empty")
		return
	}

	stat, err := os.Stat(dirFlag)

	if os.IsExist(err) {
		fmt.Printf("[*] The %q is already exist\n", dirFlag)
		return
	}

	if err == nil {
		if stat.IsDir() {
			fmt.Printf("[*] The %q is already exist\n", dirFlag)
			return
		}

		if !stat.IsDir() {
			PrintfError("Not a directory: %q\n", dirFlag)
			return
		}
	}

	_, err = getman.NewClient(dirFlag)

	if err != nil {
		PrintfError("NewClient: %s\n", err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
