package apply

import (
	"fmt"
	"sort"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"github.com/ep0ll/rebaze/internal/remote"
)

// Result is the outcome of applying a plan.
type Result struct {
	Image       v1.Image
	Source      string
	Destination string
	Steps       []string
}

func Apply(plan *Plan) (*Result, error) {
	if plan == nil {
		return nil, fmt.Errorf("plan is nil")
	}
	img, _, err := remote.Image(plan.Image)
	if err != nil {
		return nil, err
	}

	res := &Result{Source: plan.Image, Destination: plan.Tag}
	if res.Destination == "" {
		res.Destination = plan.Image
	}

	img, steps, err := applyLayers(img, plan.Layers)
	if err != nil {
		return nil, err
	}
	res.Steps = append(res.Steps, steps...)

	img, cfgSteps, err := applyConfig(img, plan.Config, plan.Annotations)
	if err != nil {
		return nil, err
	}
	res.Steps = append(res.Steps, cfgSteps...)
	res.Image = img
	return res, nil
}

func applyLayers(img v1.Image, spec *LayerMut) (v1.Image, []string, error) {
	if spec == nil {
		return img, nil, nil
	}
	var steps []string
	layers, err := img.Layers()
	if err != nil {
		return nil, nil, err
	}

	rebuild := false
	if len(spec.Delete) > 0 {
		deleteSet := map[int]struct{}{}
		for _, i := range spec.Delete {
			if i < 0 || i >= len(layers) {
				return nil, nil, fmt.Errorf("delete layer index %d out of range (0-%d)", i, len(layers)-1)
			}
			deleteSet[i] = struct{}{}
		}
		kept := make([]v1.Layer, 0, len(layers))
		for i, l := range layers {
			if _, drop := deleteSet[i]; drop {
				steps = append(steps, fmt.Sprintf("delete layer %d", i))
				continue
			}
			kept = append(kept, l)
		}
		layers = kept
		rebuild = true
	}

	if len(spec.Insert) > 0 {
		sorted := append([]LayerInsert(nil), spec.Insert...)
		sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Index < sorted[j].Index })
		for _, ins := range sorted {
			donor, _, err := remote.Image(ins.From)
			if err != nil {
				return nil, nil, err
			}
			donorLayers, err := donor.Layers()
			if err != nil {
				return nil, nil, err
			}
			if len(donorLayers) == 0 {
				return nil, nil, fmt.Errorf("insert source %s has no layers", ins.From)
			}
			idx := ins.Index
			if idx < 0 || idx > len(layers) {
				return nil, nil, fmt.Errorf("insert index %d out of range (0-%d)", idx, len(layers))
			}
			next := make([]v1.Layer, 0, len(layers)+len(donorLayers))
			next = append(next, layers[:idx]...)
			next = append(next, donorLayers...)
			next = append(next, layers[idx:]...)
			layers = next
			steps = append(steps, fmt.Sprintf("insert %d layer(s) from %s at index %d", len(donorLayers), ins.From, idx))
		}
		rebuild = true
	}

	if rebuild {
		img, err = replaceLayers(img, layers)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, src := range spec.Append {
		donor, _, err := remote.Image(src)
		if err != nil {
			return nil, nil, err
		}
		donorLayers, err := donor.Layers()
		if err != nil {
			return nil, nil, err
		}
		img, err = mutate.AppendLayers(img, donorLayers...)
		if err != nil {
			return nil, nil, fmt.Errorf("append layers from %s: %w", src, err)
		}
		steps = append(steps, fmt.Sprintf("append %d layer(s) from %s", len(donorLayers), src))
	}

	return img, steps, nil
}

func replaceLayers(img v1.Image, layers []v1.Layer) (v1.Image, error) {
	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}
	out := empty.Image
	out, err = mutate.ConfigFile(out, cfg)
	if err != nil {
		return nil, fmt.Errorf("preserve config while replacing layers: %w", err)
	}
	if len(layers) == 0 {
		return out, nil
	}
	out, err = mutate.AppendLayers(out, layers...)
	if err != nil {
		return nil, fmt.Errorf("rebuild layers: %w", err)
	}
	return out, nil
}
