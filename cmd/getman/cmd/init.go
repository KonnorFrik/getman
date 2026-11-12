/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	getman "github.com/KonnorFrik/getman/client"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a base directory with name 'dir'",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run:   _InitCmd,
}

func _InitCmd(cmd *cobra.Command, args []string) {
	if flagDirectory == "" {
		PrintfCobraError(cmd, "Flag 'dir' cannot be empty\n")
		return
	}

	stat, err := os.Stat(flagDirectory)
	if err == nil {
		if !stat.IsDir() {
			PrintfError("Not a directory: %q\n", flagDirectory)
			return
		}

		fmt.Printf("[*] The %q is already exist\n", flagDirectory)
		return
	}

	if !os.IsNotExist(err) {
		PrintfError("Can't stat %q: %s\n", flagDirectory, err)
		return
	}

	if _, err := getman.NewClient(flagDirectory); err != nil {
		PrintfError("NewClient: %s\n", err)
		return
	}

	fmt.Printf("Initialized %q\n", flagDirectory)
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
