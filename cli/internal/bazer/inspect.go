package bazer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:     "inspect <reference>",
	Short:   "Inspect manifest or index information for a container reference",
	Example: "  rebaze inspect ubuntu:latest\n  rebaze inspect ghcr.io/example/app@sha256:deadbeef",
	Args:    cobra.ExactArgs(1),
	RunE:    runInspect,
}

func runInspect(cmd *cobra.Command, args []string) error {
	ref, err := name.ParseReference(args[0], name.WeakValidation)
	if err != nil {
		return fmt.Errorf("parse reference: %w", err)
	}

	desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("fetch descriptor: %w", err)
	}

	var payload any
	switch {
	case desc.MediaType.IsIndex():
		idx, err := desc.ImageIndex()
		if err != nil {
			return fmt.Errorf("load index: %w", err)
		}
		mf, err := idx.IndexManifest()
		if err != nil {
			return fmt.Errorf("index manifest: %w", err)
		}
		payload = mf
	case desc.MediaType.IsImage():
		img, err := desc.Image()
		if err != nil {
			return fmt.Errorf("load image: %w", err)
		}
		mf, err := img.Manifest()
		if err != nil {
			return fmt.Errorf("image manifest: %w", err)
		}
		payload = mf
	default:
		return errors.New("unsupported media type: " + string(desc.MediaType))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
