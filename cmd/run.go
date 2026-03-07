package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/vikingpingvin/chainrun/internal/types"
)

var runCmd = &cobra.Command{
	Use:   "run [workflow-name]",
	Short: "Execute a workflow once, or all workflows if no name given",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var names []string
		if len(args) == 1 {
			names = []string{args[0]}
		} else {
			names = eng.Workflows()
		}

		event := types.TriggerEvent{Type: "manual", FiredAt: time.Now()}
		for _, name := range names {
			runCtx, err := eng.RunOnce(cmd.Context(), name, event)
			if err != nil {
				return err
			}
			for _, step := range runCtx.Steps {
				cmd.Printf("Step %q: status=%s stdout=%q stderr=%q\n",
					step.StepID, step.Status, step.Stdout, step.Stderr)
			}
		}
		return nil
	},
}
