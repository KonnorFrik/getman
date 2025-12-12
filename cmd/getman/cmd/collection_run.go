/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	getman "github.com/KonnorFrik/getman/client"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <collection_name> ...<collection_name>",
	Short: "Run a given collection of requests",
	Long:  `Run a given collection of requests from current dir`,
	Args:  cobra.MinimumNArgs(1),
	Run:   _RunCmd,
}

// var flagCollectionRunWithEnvName string

func _RunCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	// var envNames []string
	// if flagCollectionRunWithEnvName != "" {
	// 	for name := range strings.SplitSeq(flagCollectionRunWithEnvName, ",") {
	// 		name = strings.TrimSpace(name)
	// 		envNames = append(envNames, name)
	// 	}
	// }

	for ind, collectionName := range args {
		// if ind < len(envNames) {
		// }

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
	// runCmd.Flags().StringVarP(&flagCollectionRunWithEnvName, "env", "e", "", "Load specified collection and run collections. Also allow many names: \"e1,e2,e3\". If names less than collections - run use env from collection.")
}
