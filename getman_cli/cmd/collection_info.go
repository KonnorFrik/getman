/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// collectionInfoCmd represents the info command
var collectionInfoCmd = &cobra.Command{
	Use:   "info <collection> ...<collection>",
	Short: "Show info about collection",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: _CollectionInfoCmd,
}

func _CollectionInfoCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	methodCounter := make(map[string]int)

	for _, name := range args {
		coll, err := client.LoadCollection(name)

		if err != nil {
			PrintfError("Can't load collection %q: %s\n", name, err)
			continue
		}

		for _, req := range coll.Items {
			methodCounter[req.Request.Method]++
		}

		fmt.Printf("%s:\n", coll.Name)
		fmt.Printf("\tDescription: %s\n", coll.Description)
		fmt.Printf("\tlinked env: %s\n", coll.EnvName)
		fmt.Printf("\trequests: %d\n", len(coll.Items))

		for method, count := range methodCounter {
			fmt.Printf("\t\t%s: %d\n", method, count)
		}

		println()
	}
}

func init() {
	rootCmd.AddCommand(collectionInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
