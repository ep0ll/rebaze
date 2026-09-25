package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a structured MutationBundle / patch to an image",
	Long:  "Apply executes a MutationBundle (see schema.json) against a target image or index.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("apply: not yet implemented — see schema.json and internal/apply for the planned design")
	},
}

var apply = applyCmd
