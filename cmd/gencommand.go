package cmd

import (
	"github.com/spf13/cobra"
)

var gencommandCmd = &cobra.Command{
	Use:   "gencommand",
	Short: "Generate commands for talosctl.",
}

func init() {
	rootCmd.AddCommand(gencommandCmd)
}
