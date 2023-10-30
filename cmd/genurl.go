package cmd

import (
	"github.com/spf13/cobra"
)

var genurlCmd = &cobra.Command{
	Use:   "genurl",
	Short: "Generate URL for Talos installer or ISO",
}

func init() {
	rootCmd.AddCommand(genurlCmd)
}
