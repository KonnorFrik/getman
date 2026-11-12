/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// placeholderCmd represents the placeholder command.
var placeholderCmd = &cobra.Command{
	Use:   "placeholder",
	Short: "Work with path placeholders",
	Long: `Path placeholders are key-value pairs used to replace segments like ":id"
in request URLs, for example: https://api.example.com/users/:id`,
}

func init() {
	rootCmd.AddCommand(placeholderCmd)
}
