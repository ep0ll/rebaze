package apply

import (
	"fmt"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

func applyConfig(img v1.Image, spec *ConfigMut, anns map[string]string) (v1.Image, []string, error) {
	var steps []string
	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, nil, fmt.Errorf("read config: %w", err)
	}
	if cfg == nil {
		return nil, nil, fmt.Errorf("image config is nil")
	}

	changed := false
	if spec != nil {
		if spec.SetEnv != nil {
			cfg.Config.Env = mergeEnv(cfg.Config.Env, spec.SetEnv)
			steps = append(steps, fmt.Sprintf("setEnv %d key(s)", len(spec.SetEnv)))
			changed = true
		}
		if len(spec.UnsetEnv) > 0 {
			cfg.Config.Env = unsetEnv(cfg.Config.Env, spec.UnsetEnv)
			steps = append(steps, fmt.Sprintf("unsetEnv %v", spec.UnsetEnv))
			changed = true
		}
		if spec.SetLabel != nil {
			if cfg.Config.Labels == nil {
				cfg.Config.Labels = map[string]string{}
			}
			for k, v := range spec.SetLabel {
				cfg.Config.Labels[k] = v
			}
			steps = append(steps, fmt.Sprintf("setLabel %d key(s)", len(spec.SetLabel)))
			changed = true
		}
		if len(spec.UnsetLabel) > 0 {
			for _, k := range spec.UnsetLabel {
				delete(cfg.Config.Labels, k)
			}
			steps = append(steps, fmt.Sprintf("unsetLabel %v", spec.UnsetLabel))
			changed = true
		}
		if spec.SetEntrypoint != nil {
			cfg.Config.Entrypoint = append([]string(nil), spec.SetEntrypoint...)
			steps = append(steps, "setEntrypoint")
			changed = true
		}
		if spec.SetCmd != nil {
			cfg.Config.Cmd = append([]string(nil), spec.SetCmd...)
			steps = append(steps, "setCmd")
			changed = true
		}
		if spec.SetUser != "" {
			cfg.Config.User = spec.SetUser
			steps = append(steps, "setUser "+spec.SetUser)
			changed = true
		}
		if spec.SetWorkdir != "" {
			cfg.Config.WorkingDir = spec.SetWorkdir
			steps = append(steps, "setWorkdir "+spec.SetWorkdir)
			changed = true
		}
	}

	if changed {
		img, err = mutate.ConfigFile(img, cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("write config: %w", err)
		}
	}

	if len(anns) > 0 {
		img = mutate.Annotations(img, anns).(v1.Image)
		steps = append(steps, fmt.Sprintf("annotations %d key(s)", len(anns)))
	}

	return img, steps, nil
}

func mergeEnv(existing []string, set map[string]string) []string {
	idx := map[string]int{}
	out := append([]string(nil), existing...)
	for i, kv := range out {
		k, _, ok := strings.Cut(kv, "=")
		if ok {
			idx[k] = i
		}
	}
	for k, v := range set {
		entry := k + "=" + v
		if i, ok := idx[k]; ok {
			out[i] = entry
			continue
		}
		out = append(out, entry)
	}
	return out
}

func unsetEnv(existing []string, keys []string) []string {
	drop := map[string]struct{}{}
	for _, k := range keys {
		drop[k] = struct{}{}
	}
	out := make([]string, 0, len(existing))
	for _, kv := range existing {
		k, _, _ := strings.Cut(kv, "=")
		if _, skip := drop[k]; skip {
			continue
		}
		out = append(out, kv)
	}
	return out
}
