package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show mutation / audit history for an image or bundle",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("history: not yet implemented")
	},
}

var history = historyCmd
