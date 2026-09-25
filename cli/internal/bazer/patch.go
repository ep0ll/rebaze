package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Create a patch artifact describing mutations",
	Long:  "Create a verifiable PatchSpec / MutationBundle that can later be applied or audited.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("patch: not yet implemented")
	},
}

var patch = patchCmd
