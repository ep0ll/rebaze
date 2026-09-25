package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze/internal/history"
	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:     "copy <src> <dst>",
	Short:   "Copy an image between references / registries",
	Example: "  rebaze copy alpine:latest localhost:5000/alpine:latest",
	Args:    cobra.ExactArgs(2),
	RunE:    runCopy,
}

func runCopy(cmd *cobra.Command, args []string) error {
	ref, digest, err := remote.Copy(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s@%s\n", ref.Context().Name(), digest)
	path, err := history.DefaultPath()
	if err == nil {
		_ = history.Append(path, history.Event{
			Action:      "copy",
			Source:      args[0],
			Destination: args[1],
			Digest:      digest.String(),
		})
	}
	return nil
}
