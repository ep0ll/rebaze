package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Efficiently copy an image between registries",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("copy: not yet implemented")
	},
}

var copy = copyCmd
