package remote

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func Auth() remote.Option {
	return remote.WithAuthFromKeychain(authn.DefaultKeychain)
}

func Parse(ref string) (name.Reference, error) {
	parsed, err := name.ParseReference(ref, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("parse reference %q: %w", ref, err)
	}
	return parsed, nil
}

func Image(ref string) (v1.Image, name.Reference, error) {
	parsed, err := Parse(ref)
	if err != nil {
		return nil, nil, err
	}
	img, err := remote.Image(parsed, Auth())
	if err != nil {
		return nil, nil, fmt.Errorf("fetch image %s: %w", parsed, err)
	}
	return img, parsed, nil
}

func Write(ref string, img v1.Image) (name.Reference, v1.Hash, error) {
	parsed, err := Parse(ref)
	if err != nil {
		return nil, v1.Hash{}, err
	}
	if err := remote.Write(parsed, img, Auth()); err != nil {
		return nil, v1.Hash{}, fmt.Errorf("push %s: %w", parsed, err)
	}
	digest, err := img.Digest()
	if err != nil {
		return nil, v1.Hash{}, err
	}
	return parsed, digest, nil
}

func Copy(src, dst string) (name.Reference, v1.Hash, error) {
	img, _, err := Image(src)
	if err != nil {
		return nil, v1.Hash{}, err
	}
	return Write(dst, img)
}
