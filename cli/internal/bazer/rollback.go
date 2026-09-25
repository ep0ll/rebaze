package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze/internal/history"
	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var rollbackRef string

var rollbackCmd = &cobra.Command{
	Use:   "rollback <reference>",
	Short: "Roll a tag back to the previous digest recorded in local history",
	Args:  cobra.ExactArgs(1),
	RunE:  runRollback,
}

func runRollback(cmd *cobra.Command, args []string) error {
	rollbackRef = args[0]
	path, err := history.DefaultPath()
	if err != nil {
		return err
	}
	ev, err := history.LastFor(path, rollbackRef)
	if err != nil {
		return err
	}
	if ev.Previous == "" {
		return fmt.Errorf("history entry for %s has no previous digest to restore", rollbackRef)
	}

	src := ev.Destination + "@" + ev.Previous
	if ev.Destination == "" {
		src = ev.Source + "@" + ev.Previous
	}
	img, _, err := remote.Image(src)
	if err != nil {
		// previous may already be a full digest reference
		img, _, err = remote.Image(ev.Previous)
		if err != nil {
			return fmt.Errorf("fetch previous digest %s: %w", ev.Previous, err)
		}
	}
	ref, digest, err := remote.Write(rollbackRef, img)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rolled back %s -> %s@%s\n", rollbackRef, ref.Context().Name(), digest)
	return history.Append(path, history.Event{
		Action:      "rollback",
		Source:      ev.Digest,
		Destination: rollbackRef,
		Digest:      digest.String(),
		Previous:    ev.Digest,
	})
}

var rollback = rollbackCmd
