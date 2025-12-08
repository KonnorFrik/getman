/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// envUnsetCmd represents the unset command
var envUnsetCmd = &cobra.Command{
	Use:   "unset <env> <key> ...<key>",
	Short: "Delete a pair key-value from environment",
	Long: ``,
	Args: cobra.MinimumNArgs(2),
	Run: _EnvUnsetCmd,
}

func _EnvUnsetCmd(cmd *cobra.Command, args []string) {
	if dirFlag == "" {
		PrintfCobraError(cmd, "Flag 'dir' cannot be empty")
		return
	}

	pathStat, err := os.Stat(dirFlag)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	if !pathStat.IsDir() {
		PrintfError("not a directory: %s\n", dirFlag)
		return
	}

	client, err := getman.NewClient(dirFlag)

	if err != nil {
		PrintfError("NewClient: %s\n", err)
		return
	}

	envName := args[0]

	err = client.LoadLocalEnvironment(envName)

	if err != nil {
		PrintfError("Can't load environment %q - %s\n", envName, err)
		return
	}

	env := client.GetCurrentEnvironment()
	var deleted int

	for _, key := range args[1:] {
		_, exist := env.Variables[key]

		if !exist {
			continue
		}

		delete(env.Variables, key)
		deleted++
	}

	err = client.SaveEnvironment(env)

	if err != nil {
		PrintfError("Can't save environment %q - %s\n", env, err)
		return
	}

	fmt.Printf("From %s: Deleted: %d\n", env.Name, deleted)
}

func init() {
	envCmd.AddCommand(envUnsetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unsetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
