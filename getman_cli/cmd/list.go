/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

/*
list all collections from storage.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all collection names.",
	Long: `Print all names of collections stored in getman directory.`,
	Args: cobra.NoArgs,
	Run: _ListCmd,
}

func _ListCmd(cmd *cobra.Command, args []string) {
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

	collectionNames, err := client.ListCollections()

	if err != nil {
		PrintfError("ListCollections: %s\n", err)
		return
	}

	for ind, name := range collectionNames {
		fmt.Printf("%d: %s\n", ind, name)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// listCmd.Flags().StringVar(&dirFlag, "dir", ".getman", "Use specific getman directory.")
}
