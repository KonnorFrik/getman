/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link <collection_name> <env_name>",
	Short: "Set new env for collection",
	Long: ``,
	Args: cobra.ExactArgs(2),
	Run: _CollectionLinkCmd,
}

func _CollectionLinkCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	var (
		collectionName = args[0]
		environmentName = args[1]
	)

	err = client.LoadLocalEnvironment(environmentName)

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

	collection, err := client.LoadCollection(collectionName)

	if err != nil {
		PrintfError("Can't load collection: %s\n", err)
		return
	}

	collection.EnvName = environmentName
	err = client.SaveCollection(collection)

	if err != nil {
		PrintfError("Can't save collection: %s\n", err)
		return
	}
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
