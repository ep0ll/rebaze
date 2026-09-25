package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Roll back a previous mutation (pointer-revert or inverse-apply)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("rollback: not yet implemented")
	},
}

var rollback = rollbackCmd
