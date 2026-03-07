package cmd

import (
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start all workflow triggers and run continuously",
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
	},
}
