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
	fmt.Println("run called")
	fmt.Printf("args: %v\n", args)
	fmt.Printf("use dir: %s\n", dirFlag)
	pathStat, err := os.Stat(dirFlag)

	if err != nil {
		PrintfCobraError(cmd, "%s\n", err)
		return
	}

	if !pathStat.IsDir() {
		PrintfCobraError(cmd, "not a directory: %s\n", dirFlag)
		return
	}

	client, err := getman.NewClient(dirFlag)

	if err != nil {
		PrintfCobraError(cmd, "NewClient: %s\n", err)
		return
	}

	for ind, collectionName := range args {
		result, err := client.ExecuteCollection(collectionName)

		if err != nil {
			PrintfError("[!] #%d - %s\n", ind, err)
			continue
		}

		getman.PrintExecutionResult(result)
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
	runCmd.Flags().StringVar(&dirFlag, "dir", ".getman", "Use specific getman directory.")
}
