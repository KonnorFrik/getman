/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var reqListCmd = &cobra.Command{
	Use:   "list <collection> ...<collection>",
	Short: "List all requests from collection",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: _ReqListCmd,
}

func _ReqListCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	for _, name := range args {
		colleciton, err := client.LoadCollection(name)

		if err != nil {
			PrintfError("can't load collection %q - %s\n", name, err)
			continue
		}

		fmt.Printf("%s:\n", name)

		if len(colleciton.Items) == 0 {
			fmt.Printf("\tEmpty\n")
			continue
		}

		for ind, req := range colleciton.Items {
			fmt.Printf("\t%d: %s -> %s %s\n", ind, req.Name, req.Request.Method, req.Request.URL)
		}

		println()
	}
}

func init() {
	reqCmd.AddCommand(reqListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
