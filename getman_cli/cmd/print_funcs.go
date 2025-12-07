package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// PrintfCobraError - Print formated message always into Stderr with prefix from cobra lib.
// Print usage at the end
func PrintfCobraError(cmd *cobra.Command, msg string, args ...any) {
	fmt.Fprintf(os.Stderr, cmd.ErrPrefix() + " " + msg, args...)
	cmd.Usage()
}

// PrintfError - Print formated message always into Stderr with prefix.
func PrintfError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
}
