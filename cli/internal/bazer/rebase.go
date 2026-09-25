package bazer

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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

This is the same fundamental operation performed by crane rebase and the
Cloud Native Buildpacks rebaser. It is extremely efficient when the registry
supports cross-repository blob mounts.

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
	origRef, err := name.ParseReference(args[0], name.WeakValidation)
	if err != nil {
		return fmt.Errorf("parse original image: %w", err)
	}

	newBaseRef, err := name.ParseReference(rebaseNewBase, name.WeakValidation)
	if err != nil {
		return fmt.Errorf("parse new-base: %w", err)
	}

	authOpt := remote.WithAuthFromKeychain(authn.DefaultKeychain)

	origImg, err := remote.Image(origRef, authOpt)
	if err != nil {
		return fmt.Errorf("fetch original image: %w", err)
	}

	var oldBaseImg v1.Image
	if rebaseOldBase != "" {
		oldBaseRef, err := name.ParseReference(rebaseOldBase, name.WeakValidation)
		if err != nil {
			return fmt.Errorf("parse old-base: %w", err)
		}
		oldBaseImg, err = remote.Image(oldBaseRef, authOpt)
		if err != nil {
			return fmt.Errorf("fetch old-base image: %w", err)
		}
	} else {
		// Attempt to use OCI base annotations when present
		mf, err := origImg.Manifest()
		if err != nil {
			return fmt.Errorf("read original manifest: %w", err)
		}
		baseDigest, ok1 := mf.Annotations["org.opencontainers.image.base.digest"]
		baseName, ok2 := mf.Annotations["org.opencontainers.image.base.name"]
		if !ok1 || !ok2 {
			return fmt.Errorf("--old-base is required when the image lacks OCI base annotations")
		}
		oldBaseRef, err := name.ParseReference(fmt.Sprintf("%s@%s", baseName, baseDigest), name.WeakValidation)
		if err != nil {
			return fmt.Errorf("parse annotated base: %w", err)
		}
		oldBaseImg, err = remote.Image(oldBaseRef, authOpt)
		if err != nil {
			return fmt.Errorf("fetch annotated old-base: %w", err)
		}
	}

	newBaseImg, err := remote.Image(newBaseRef, authOpt)
	if err != nil {
		return fmt.Errorf("fetch new-base image: %w", err)
	}

	rebased, err := mutate.Rebase(origImg, oldBaseImg, newBaseImg)
	if err != nil {
		return fmt.Errorf("rebase: %w", err)
	}

	target := origRef
	if rebaseTag != "" {
		target, err = name.ParseReference(rebaseTag, name.WeakValidation)
		if err != nil {
			return fmt.Errorf("parse --tag: %w", err)
		}
	}

	if err := remote.Write(target, rebased, authOpt); err != nil {
		return fmt.Errorf("push rebased image: %w", err)
	}

	digest, err := rebased.Digest()
	if err != nil {
		return fmt.Errorf("compute digest: %w", err)
	}

	fmt.Printf("%s@%s\n", target.Context().Name(), digest)
	return nil
}

// Keep the old helper name for residual references.
var rebase = rebaseCmd
