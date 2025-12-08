/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all environments",
	Long: ``,
	Run: _EnvListCmd,
}

func _EnvListCmd(cmd *cobra.Command, args []string) {
	if dirFlag == "" {
		PrintfCobraError(cmd, "Flag 'dir' cannot be empty")
		return
	}

	pathStat, err := os.Stat(dirFlag)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	if !pathStat.IsDir() {
		PrintfError("Not a directory: %s\n", dirFlag)
		return
	}

	client, err := getman.NewClient(dirFlag)

	if err != nil {
		PrintfError("NewClient: %s\n", err)
		return
	}

	envNames, err := client.ListEnvironments()

	if err != nil {
		PrintfError("ListCollections: %s\n", err)
		return
	}

	for ind, name := range envNames {
		fmt.Printf("%d: %s\n", ind, name)
	}
}

func init() {
	envCmd.AddCommand(envListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
