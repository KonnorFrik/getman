/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package cmd

import (
	getman "github.com/KonnorFrik/getman/client"
	"github.com/spf13/cobra"
)

// envInitCmd represents the init command
var envInitCmd = &cobra.Command{
	Use:   "init <name> ...<name>",
	Short: "Create a new environment",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run:   _EnvInitCmd,
}

func _EnvInitCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	for _, name := range args {
		err := createEnvInStorage(client, name)

		if err != nil {
			PrintfError("Can't create env %q - %s\n", name, err)
		}
	}
}

func init() {
	envCmd.AddCommand(envInitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createEnvInStorage(cl *getman.Client, envName string) error {
	env := &getman.Environment{
		Name: envName,
	}

	err := cl.SaveEnvironment(env)

	return err
}
