/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// PrintfCobraError - Print formated message always into Stderr with prefix from cobra lib.
// Print usage at the end
func PrintfCobraError(cmd *cobra.Command, format string, args ...any) {
	fmt.Fprintf(os.Stderr, cmd.ErrPrefix() + " " + format, args...)
	cmd.Usage()
}

// FormatCobraError - create a error from 'msg' as format. Add prefix and cobra usage.
func CobraErrorf(cmd *cobra.Command, msg string, args ...any) error {
	format := fmt.Sprintf("%s %s%s", cmd.ErrPrefix(), msg, cmd.UsageString())
	return fmt.Errorf(format, args)
}

// PrintfError - Print formated message always into Stderr with prefix.
func PrintfError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
