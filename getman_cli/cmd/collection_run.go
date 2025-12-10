/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (

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

func _RunCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
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
