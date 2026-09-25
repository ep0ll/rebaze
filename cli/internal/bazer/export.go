package bazer

import (
	"encoding/json"
	"fmt"

	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:     "export <reference>",
	Short:   "Export image config, layers, and digest as JSON",
	Example: "  rebaze export alpine:latest",
	Args:    cobra.ExactArgs(1),
	RunE:    runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	img, ref, err := remote.Image(args[0])
	if err != nil {
		return err
	}
	cfg, err := img.ConfigFile()
	if err != nil {
		return err
	}
	mf, err := img.Manifest()
	if err != nil {
		return err
	}
	digest, err := img.Digest()
	if err != nil {
		return err
	}
	size, err := img.Size()
	if err != nil {
		return err
	}

	out := map[string]any{
		"reference": ref.String(),
		"digest":    digest.String(),
		"size":      size,
		"manifest":  mf,
		"config":    cfg,
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return err
	}
	fmt.Fprintln(cmd.ErrOrStderr(), digest.String())
	return nil
}

var export = exportCmd
