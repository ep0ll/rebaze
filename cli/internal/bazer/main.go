package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rebaze",
	Short: "Efficient OCI / Docker image mutation and rebasing",
	Long: `rebaze is a CLI for efficient OCI image mutation and rebasing.

It lets you replace base layers or apply structured patches without a full rebuild,
making it ideal for rolling out base-image security fixes and OS updates across
many application images.

Core commands:
  inspect   Inspect manifest / index information
  rebase    Rebase an image onto a new base image

Additional commands (apply, patch, preview, history, rollback, sign, …)
are planned and currently return a clear "not implemented" status.`,
	Version: fmt.Sprintf("%s (%s)", rebaze.Version, rebaze.Revision),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.SetVersionTemplate("rebaze {{.Version}}\n")

	rootCmd.AddCommand(
		inspectCmd,
		rebaseCmd,
		applyCmd,
		patchCmd,
		previewCmd,
		historyCmd,
		rollbackCmd,
		signCmd,
		copyCmd,
		exportCmd,
	)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
