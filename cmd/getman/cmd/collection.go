/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import "github.com/spf13/cobra"

// collectionCmd represents the collection command.
var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Work with request collections",
}

func init() {
	rootCmd.AddCommand(collectionCmd)
}
