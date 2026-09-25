package bazer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ep0ll/rebaze/internal/apply"
	"github.com/spf13/cobra"
)

var (
	patchImage string
	patchOut   string
	patchEnv   []string
	patchUser  string
	patchLabel []string
)

var patchCmd = &cobra.Command{
	Use:   "patch --image <ref>",
	Short: "Create a mutation plan JSON from flags",
	Long:  "Generate a plan file that apply / preview can consume.",
	RunE:  runPatch,
}

func init() {
	patchCmd.Flags().StringVar(&patchImage, "image", "", "Source image reference")
	patchCmd.Flags().StringVar(&patchOut, "out", "plan.json", "Output plan path")
	patchCmd.Flags().StringArrayVar(&patchEnv, "set-env", nil, "KEY=VALUE environment mutations")
	patchCmd.Flags().StringVar(&patchUser, "set-user", "", "Config user")
	patchCmd.Flags().StringArrayVar(&patchLabel, "set-label", nil, "KEY=VALUE label mutations")
	_ = patchCmd.MarkFlagRequired("image")
}

func runPatch(cmd *cobra.Command, args []string) error {
	plan := &apply.Plan{
		SchemaVersion: "1.0",
		Kind:          "MutationPlan",
		Image:         patchImage,
		Config:        &apply.ConfigMut{},
	}
	if len(patchEnv) > 0 {
		plan.Config.SetEnv = map[string]string{}
		for _, kv := range patchEnv {
			k, v, ok := splitKV(kv)
			if !ok {
				return fmt.Errorf("invalid --set-env %q (want KEY=VALUE)", kv)
			}
			plan.Config.SetEnv[k] = v
		}
	}
	if len(patchLabel) > 0 {
		plan.Config.SetLabel = map[string]string{}
		for _, kv := range patchLabel {
			k, v, ok := splitKV(kv)
			if !ok {
				return fmt.Errorf("invalid --set-label %q (want KEY=VALUE)", kv)
			}
			plan.Config.SetLabel[k] = v
		}
	}
	plan.Config.SetUser = patchUser

	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(patchOut, append(data, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wrote %s at %s\n", patchOut, time.Now().UTC().Format(time.RFC3339))
	return nil
}

func splitKV(kv string) (string, string, bool) {
	for i := 0; i < len(kv); i++ {
		if kv[i] == '=' {
			return kv[:i], kv[i+1:], true
		}
	}
	return "", "", false
}

var patch = patchCmd
