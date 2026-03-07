package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vikingpingvin/chainrun/config"
	"github.com/vikingpingvin/chainrun/engine"
)

var (
	eng    engine.Engine
	loader config.Loader
)

var rootCmd = &cobra.Command{
	Use:   "chainrun",
	Short: "ChainRun — workflow automation engine",
}

// Execute wires the engine and loader into the CLI and runs the root command.
func Execute(e engine.Engine, l config.Loader) error {
	eng = e
	loader = l

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)

	return rootCmd.Execute()
}
