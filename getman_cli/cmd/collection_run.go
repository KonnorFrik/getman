/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run ...",
	Short: "Run a given collection of requests",
	Long: `Run a given collection of requests from current dir`,
	Args: cobra.MinimumNArgs(1),
	Run: _RunCmd,
}

var dirFlag string

func _RunCmd(cmd *cobra.Command, args []string) {
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
		PrintfError("not a directory: %s\n", dirFlag)
		return
	}

	client, err := getman.NewClient(dirFlag)

	if err != nil {
		PrintfError("NewClient: %s\n", err)
		return
	}

	for ind, collectionName := range args {
		result, err := client.ExecuteCollection(collectionName)

		if err != nil {
			PrintfError("[!] #%d - %s\n", ind, err)
			continue
		}

		getman.PrintExecutionResult(result)
		println()
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
