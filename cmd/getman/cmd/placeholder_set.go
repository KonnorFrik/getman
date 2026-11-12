/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var placeholderSetCmd = &cobra.Command{
	Use:   "set <key=value> ...<key=value>",
	Short: "Set path placeholder values",
	Args:  cobra.MinimumNArgs(1),
	Run:   _PlaceholderSetCmd,
}

func _PlaceholderSetCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)
	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	updated := 0

	for _, pair := range args {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			PrintfError("Invalid pair %q, expected key=value\n", pair)
			continue
		}

		if err := client.SetPlaceholder(parts[0], parts[1]); err != nil {
			PrintfError("Can't save placeholder %q: %s\n", parts[0], err)
			continue
		}

		updated++
	}

	fmt.Printf("Updated: %d\n", updated)
}

func init() {
	placeholderCmd.AddCommand(placeholderSetCmd)
}
