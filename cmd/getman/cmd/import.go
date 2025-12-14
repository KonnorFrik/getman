/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <path> ...<path>",
	Short: "Import a collection from another format",
	Long: `Import collection from <path>.
If a collection with the same name already exists, it will be overwritten !
Can import from:
- Postman json`,
	Args: cobra.MinimumNArgs(1),
	Run: _ImportCmd,
}

var (
	flagImportFromPostman bool
)

func _ImportCmd(cmd *cobra.Command, args []string) {
	if !flagImportFromPostman {
		PrintfCobraError(cmd, "Can't determine format for import form. Need one of flags to be used.")
		return
	}

	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	var imported int

	for _, path := range args {
		coll, err := client.ImportFromPostman(path)

		if err != nil {
			PrintfError("Can't import from %s: %s\n", path, err)

			if flagExitOnError {
				fmt.Printf("Exit on error.")
				return
			}

			continue
		}

		err = client.SaveCollection(coll)

		if err != nil {
			PrintfError("Can't save collection: %s\n", err)

			if flagExitOnError {
				fmt.Printf("Exit on error.")
				return
			}
		}

		imported++
	}

	fmt.Printf("Imported: %d collections\n", imported)
}

// TODO: for more, than one flag - need check for: only one flag enaled
// TODO: if collection already imported and exist, flag --overwrite for update, default, do nothing ?

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&flagImportFromPostman, "postman", false, "Import collection from postman json file.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

