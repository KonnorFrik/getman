/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create new collections.",
	Long: `Create a new collection in storage.`,
	Args: cobra.ExactArgs(1),
	Run: _CollectionNewCmd,
}

func _CollectionNewCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	if flagCollectionEnvName != "" {
		err = client.LoadLocalEnvironment(flagCollectionEnvName)

		if err != nil {
			if flagExitOnError {
				PrintfError("Can't find environment: %s\n", flagCollectionEnvName)
				PrintfError("Exit on error\n")
				return
			}

			err = createEnvInStorage(client, flagCollectionEnvName)

			if err != nil {
				if flagExitOnError {
					PrintfError("Can't find and create environment: %s: %s\n", flagCollectionEnvName, err)
					PrintfError("Exit on error\n")
					return
				}
			}
		}
	}

	collection := getman.Collection{
		Name: args[0],
		EnvName: flagCollectionEnvName,
		Description: flagCollectionDescription,
	}

	err = client.SaveCollection(&collection)

	if err != nil {
		PrintfError("Can't create collection: %s: %s\n", args[0], err)
	}
}

// flags
var (
	flagCollectionDescription string
	flagCollectionEnvName string
	// what if env don't exist
	// check is exist?
	// do nothing
	// flag for don't create if not exist (default - create)
)

func init() {
	// TODO: add flag for collection: description, env
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().StringVarP(&flagCollectionDescription, "description", "d", "", "Collection description.")
	newCmd.Flags().StringVarP(&flagCollectionEnvName, "link-env", "e", "", "Environment name for use.")
}
