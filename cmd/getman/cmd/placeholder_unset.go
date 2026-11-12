/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var placeholderUnsetCmd = &cobra.Command{
	Use:   "unset <key> ...<key>",
	Short: "Delete path placeholders",
	Args:  cobra.MinimumNArgs(1),
	Run:   _PlaceholderUnsetCmd,
}

func _PlaceholderUnsetCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)
	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	deleted := 0

	for _, key := range args {
		if _, ok := client.GetPlaceholder(key); !ok {
			continue
		}

		if err := client.DeletePlaceholder(key); err != nil {
			PrintfError("Can't delete placeholder %q: %s\n", key, err)
			continue
		}

		deleted++
	}

	fmt.Printf("Deleted: %d\n", deleted)
}

func init() {
	placeholderCmd.AddCommand(placeholderUnsetCmd)
}
