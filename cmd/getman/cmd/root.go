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

// global flags for any commands
var (
	// flagExitOnError if true - exit immediately.
	flagExitOnError bool
	// flagDirectory - directory path for use as storage entrypoint.
	flagDirectory string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getman",
	Short: "Getman is a CLI utility for work with http request.",
	Long: `Getman is a CLI utility for work with http request.
Features:
	- folder-storage
	- requests collections
	- variables
	- environments for variables
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.getman.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&flagDirectory, "dir", ".getman", "Use specific getman directory.")
	rootCmd.PersistentFlags().BoolVar(&flagExitOnError, "fail-fast", false, "Fast exit if any error occured.")
}

func createClientWithDirectory(cmd *cobra.Command) (*getman.Client, error) {
	if flagDirectory == "" {
		return nil, CobraErrorf(cmd, "Flag 'dir' cannot be empty")
	}

	pathStat, err := os.Stat(flagDirectory)

	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	if !pathStat.IsDir() {
		return nil, fmt.Errorf("not a directory: %s\n", flagDirectory)
	}

	client, err := getman.NewClient(flagDirectory)

	if err != nil {
		return nil, fmt.Errorf("NewClient: %s\n", err)
	}

	return client, nil
}
