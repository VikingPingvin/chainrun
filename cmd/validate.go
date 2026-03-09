package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vikingpingvin/chainrun/internal/types"
)

var validateCmd = &cobra.Command{
	Use:   "validate [workflow-name]",
	Short: "Validate a workflow config — parse, schema-check, and dry-run template rendering",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Phase 2: structural validation
		cfg, err := loader.Load(configPath)
		if err != nil {
			return err
		}
		if errs := loader.Validate(cfg); len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %s: %s\n", e.Field, e.Message)
			}
			return fmt.Errorf("%d validation error(s)", len(errs))
		}

		// Phase 3: template dry-run
		var names []string
		if len(args) == 1 {
			names = []string{args[0]}
		} else {
			names = eng.Workflows()
		}

		event := types.TriggerEvent{Type: "manual", FiredAt: time.Now()}
		for _, name := range names {
			if err := eng.DryRun(cmd.Context(), name, event); err != nil {
				return err
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Config is valid.")
		return nil
	},
}
