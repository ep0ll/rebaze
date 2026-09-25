package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Dry-run a mutation plan and show the resulting DAG",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("preview: not yet implemented")
	},
}

var preview = previewCmd
