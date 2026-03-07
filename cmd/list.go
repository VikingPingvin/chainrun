package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded workflows",
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
	},
}
