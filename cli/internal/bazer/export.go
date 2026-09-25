package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a mutation plan, SBOM, or audit ledger",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("export: not yet implemented")
	},
}

var export = exportCmd
