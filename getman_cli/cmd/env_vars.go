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

// varsCmd represents the vars command
var varsCmd = &cobra.Command{
	Use:   "vars",
	Short: "show all variables stored in specified environment",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: _EnvVarsCmd,
}

func _EnvVarsCmd(cmd *cobra.Command, args []string) {
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
		PrintfError("Not a directory: %s\n", dirFlag)
		return
	}

	client, err := getman.NewClient(dirFlag)

	if err != nil {
		PrintfError("NewClient: %s\n", err)
		return
	}

	for _, name := range args {
		err := client.LoadLocalEnvironment(name)

		if err != nil {
			PrintfError("Can't load env %q - %s\n", name, err)
			continue
		}

		env := client.GetCurrentEnvironment()

		fmt.Printf("%s:\n", name)

		if len(env.Variables) == 0 {
			fmt.Printf("\tEmpty\n")
			continue
		}

		for k, v := range env.Variables {
			fmt.Printf("\t%s: %s\n", k, v)
		}

		println()
	}

}

func init() {
	envCmd.AddCommand(varsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// varsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// varsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
