/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// envSetCmd represents the set command
var envSetCmd = &cobra.Command{
	Use:   "set <env> <key=value> ...<key=value>",
	Short: "Set pair key-value into specified environment",
	Long: `Must be at least one key-value pair`,
	Args: cobra.MinimumNArgs(2),
	Run: _EnvSetCmd,
}

func _EnvSetCmd(cmd *cobra.Command, args []string) {
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

	var (
		added int
		overwrited int
	)

	for ind, pair := range args[1:] {
		if !strings.Contains(pair, "=") {
			PrintfError("pair #%d - %s must contain '='\n", ind + 1, pair)
			continue
		}

		parts := strings.Split(pair, "=")

		if len(parts) != 2 {
			PrintfError("pair #%d - %s must contain exactly one '='\n", ind + 1, pair)
			continue
		}

		_, exist := env.Variables[parts[0]]
		env.Variables[parts[0]] = parts[1]

		if exist {
			overwrited++

		} else {
			added++
		}
	}

	err = client.SaveEnvironment(env)

	if err != nil {
		PrintfError("Can't save environment %q - %s\n", env, err)
		return
	}

	fmt.Printf("Into %s: Added %d, Overwrited: %d, Summary: %d\n", env.Name, added, overwrited, added + overwrited)
}

func init() {
	envCmd.AddCommand(envSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
