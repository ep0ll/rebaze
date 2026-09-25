package bazer

import (
	"fmt"

	hist "github.com/ep0ll/rebaze/internal/history"
	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <reference>",
	Short: "Roll a tag back to the previous digest recorded in local history",
	Args:  cobra.ExactArgs(1),
	RunE:  runRollback,
}

func runRollback(cmd *cobra.Command, args []string) error {
	refName := args[0]
	path, err := hist.DefaultPath()
	if err != nil {
		return err
	}
	ev, err := hist.LastFor(path, refName)
	if err != nil {
		return err
	}
	if ev.Previous == "" {
		return fmt.Errorf("history entry for %s has no previous digest to restore", refName)
	}

	src := ev.Destination + "@" + ev.Previous
	if ev.Destination == "" {
		src = ev.Source + "@" + ev.Previous
	}
	img, _, err := remote.Image(src)
	if err != nil {
		img, _, err = remote.Image(ev.Previous)
		if err != nil {
			return fmt.Errorf("fetch previous digest %s: %w", ev.Previous, err)
		}
	}
	ref, digest, err := remote.Write(refName, img)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rolled back %s -> %s@%s\n", refName, ref.Context().Name(), digest)
	return hist.Append(path, hist.Event{
		Action:      "rollback",
		Source:      ev.Digest,
		Destination: refName,
		Digest:      digest.String(),
		Previous:    ev.Digest,
	})
}
