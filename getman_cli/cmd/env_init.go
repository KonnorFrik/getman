/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// envInitCmd represents the init command
var envInitCmd = &cobra.Command{
	Use:   "init ...",
	Short: "create a init environment",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: _EnvInitCmd,
}

func _EnvInitCmd(cmd *cobra.Command, args []string) {
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

	for _, name := range args {
		env := &getman.Environment{
			Name: name,
		}

		err := client.SaveEnvironment(env)

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
