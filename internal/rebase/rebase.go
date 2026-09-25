package rebase

import (
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"github.com/ep0ll/rebaze/internal/remote"
)

func Images(original, oldBase, newBase v1.Image) (v1.Image, error) {
	out, err := mutate.Rebase(original, oldBase, newBase)
	if err != nil {
		return nil, fmt.Errorf("rebase: %w", err)
	}
	return out, nil
}

func Refs(originalRef, oldBaseRef, newBaseRef string) (v1.Image, error) {
	orig, _, err := remote.Image(originalRef)
	if err != nil {
		return nil, err
	}
	oldBase, _, err := remote.Image(oldBaseRef)
	if err != nil {
		return nil, err
	}
	newBase, _, err := remote.Image(newBaseRef)
	if err != nil {
		return nil, err
	}
	return Images(orig, oldBase, newBase)
}
