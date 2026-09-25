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

Commands:
  inspect   Inspect manifest / index information
  rebase    Rebase an image onto a new base image
  apply     Apply a JSON mutation plan
  preview   Dry-run a mutation plan
  patch     Generate a mutation plan from flags
  copy      Copy an image between registries
  export    Export config + manifest JSON
  history   Local mutation history
  rollback  Restore previous digest from history
  sign      Sign or verify a plan with ed25519`,
	Version:      fmt.Sprintf("%s (%s)", rebaze.Version, rebaze.Revision),
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

func Execute() error {
	return rootCmd.Execute()
}
