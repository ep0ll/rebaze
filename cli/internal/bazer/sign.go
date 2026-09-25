package bazer

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a MutationBundle or resulting image (cosign / notation)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("sign: not yet implemented")
	},
}

var sign = signCmd
