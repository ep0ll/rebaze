package bazer

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/ep0ll/rebaze/internal/history"
	"github.com/ep0ll/rebaze/internal/rebase"
	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var (
	rebaseOldBase string
	rebaseNewBase string
	rebaseTag     string
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase <image>",
	Short: "Rebase an image onto a new base image",
	Long: `Rebase rewrites an image so that the layers belonging to the old base
are replaced by the layers of the new base. Application layers above the base
are preserved.

Example:
  rebaze rebase my-app:1.2.3 \
    --old-base ubuntu:22.04 \
    --new-base ubuntu:24.04 \
    --tag my-app:1.2.3-rebased`,
	Args: cobra.ExactArgs(1),
	RunE: runRebase,
}

func init() {
	rebaseCmd.Flags().StringVar(&rebaseOldBase, "old-base", "", "Original base image reference (required unless OCI base annotations are present)")
	rebaseCmd.Flags().StringVar(&rebaseNewBase, "new-base", "", "New base image reference (required)")
	rebaseCmd.Flags().StringVar(&rebaseTag, "tag", "", "Tag (or full reference) to push the rebased image as")
	_ = rebaseCmd.MarkFlagRequired("new-base")
}

func runRebase(cmd *cobra.Command, args []string) error {
	origImg, origRef, err := remote.Image(args[0])
	if err != nil {
		return err
	}

	var oldBaseImg v1.Image
	if rebaseOldBase != "" {
		oldBaseImg, _, err = remote.Image(rebaseOldBase)
		if err != nil {
			return err
		}
	} else {
		mf, err := origImg.Manifest()
		if err != nil {
			return fmt.Errorf("read original manifest: %w", err)
		}
		baseDigest, ok1 := mf.Annotations["org.opencontainers.image.base.digest"]
		baseName, ok2 := mf.Annotations["org.opencontainers.image.base.name"]
		if !ok1 || !ok2 {
			return fmt.Errorf("--old-base is required when the image lacks OCI base annotations")
		}
		oldBaseImg, _, err = remote.Image(fmt.Sprintf("%s@%s", baseName, baseDigest))
		if err != nil {
			return err
		}
	}

	newBaseImg, _, err := remote.Image(rebaseNewBase)
	if err != nil {
		return err
	}

	rebased, err := rebase.Images(origImg, oldBaseImg, newBaseImg)
	if err != nil {
		return err
	}

	target := origRef.String()
	if rebaseTag != "" {
		if _, err := name.ParseReference(rebaseTag, name.WeakValidation); err != nil {
			return fmt.Errorf("parse --tag: %w", err)
		}
		target = rebaseTag
	}

	ref, digest, err := remote.Write(target, rebased)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s@%s\n", ref.Context().Name(), digest)

	path, herr := history.DefaultPath()
	if herr == nil {
		_ = history.Append(path, history.Event{
			Action:      "rebase",
			Source:      args[0],
			Destination: target,
			Digest:      digest.String(),
		})
	}
	return nil
}

var rebase = rebaseCmd
