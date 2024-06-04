package cmd

import (
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the correctness of talconfig or talos node config",
}

func init() {
	RootCmd.AddCommand(ValidateCmd)
}
