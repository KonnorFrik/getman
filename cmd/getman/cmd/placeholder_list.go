/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var placeholderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List path placeholders",
	Args:  cobra.NoArgs,
	Run:   _PlaceholderListCmd,
}

func _PlaceholderListCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)
	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	placeholders := client.ListPlaceholders()
	if len(placeholders) == 0 {
		fmt.Printf("No placeholders\n")
		return
	}

	keys := make([]string, 0, len(placeholders))
	for key := range placeholders {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, placeholders[key])
	}
}

func init() {
	placeholderCmd.AddCommand(placeholderListCmd)
}
