package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vikingpingvin/chainrun/config"
	"github.com/vikingpingvin/chainrun/engine"
)

var (
	eng           engine.Engine
	loader        config.Loader
	engineFactory func(string) (engine.Engine, error)
	configPath    string
)

var rootCmd = &cobra.Command{
	Use:          "chainrun",
	Short:        "ChainRun — workflow automation engine",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		built, err := engineFactory(configPath)
		if err != nil {
			return err
		}
		eng = built
		return nil
	},
}

// Execute wires the engine and loader into the CLI and runs the root command.
func Execute(factory func(string) (engine.Engine, error), l config.Loader) error {
	engineFactory = factory
	loader = l

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "chainrun.yaml",
		"path to workflow config file")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)

	return rootCmd.Execute()
}
