package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze/internal/apply"
	"github.com/ep0ll/rebaze/internal/history"
	"github.com/ep0ll/rebaze/internal/remote"
	"github.com/spf13/cobra"
)

var (
	applyPlan string
	applyTag  string
	applyDry  bool
)

var applyCmd = &cobra.Command{
	Use:   "apply --plan <file>",
	Short: "Apply a mutation plan to an image",
	Long: `Apply reads a JSON mutation plan and mutates the target image.

Plan example:
{
  "image": "my-app:1.0",
  "tag": "my-app:1.0-mutated",
  "config": { "setEnv": { "FOO": "bar" }, "setUser": "65532" },
  "layers": { "append": ["busybox:latest"] }
}`,
	RunE: runApply,
}

func init() {
	applyCmd.Flags().StringVar(&applyPlan, "plan", "", "Path to mutation plan JSON")
	applyCmd.Flags().StringVar(&applyTag, "tag", "", "Override destination tag from the plan")
	applyCmd.Flags().BoolVar(&applyDry, "dry-run", false, "Compute the mutation but do not push")
	_ = applyCmd.MarkFlagRequired("plan")
}

func runApply(cmd *cobra.Command, args []string) error {
	plan, err := apply.LoadPlan(applyPlan)
	if err != nil {
		return err
	}
	if applyTag != "" {
		plan.Tag = applyTag
	}

	res, err := apply.Apply(plan)
	if err != nil {
		return err
	}

	for _, s := range res.Steps {
		fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", s)
	}
	if len(res.Steps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "- no mutations requested")
	}

	if applyDry {
		digest, err := res.Image.Digest()
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "dry-run digest: %s\n", digest)
		return nil
	}

	prevDigest := ""
	if srcImg, _, err := remote.Image(res.Destination); err == nil {
		if d, err := srcImg.Digest(); err == nil {
			prevDigest = d.String()
		}
	}

	ref, digest, err := remote.Write(res.Destination, res.Image)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s@%s\n", ref.Context().Name(), digest)

	path, err := history.DefaultPath()
	if err == nil {
		_ = history.Append(path, history.Event{
			Action:      "apply",
			Source:      res.Source,
			Destination: res.Destination,
			Digest:      digest.String(),
			Previous:    prevDigest,
			Notes:       res.Steps,
		})
	}
	return nil
}

var apply = applyCmd
