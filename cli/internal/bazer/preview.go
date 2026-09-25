package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze/internal/apply"
	"github.com/spf13/cobra"
)

var previewPlan string

var previewCmd = &cobra.Command{
	Use:   "preview --plan <file>",
	Short: "Dry-run a mutation plan and print the resulting digest",
	RunE:  runPreview,
}

func init() {
	previewCmd.Flags().StringVar(&previewPlan, "plan", "", "Path to mutation plan JSON")
	_ = previewCmd.MarkFlagRequired("plan")
}

func runPreview(cmd *cobra.Command, args []string) error {
	plan, err := apply.LoadPlan(previewPlan)
	if err != nil {
		return err
	}
	res, err := apply.Apply(plan)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "source: %s\n", res.Source)
	fmt.Fprintf(cmd.OutOrStdout(), "destination: %s\n", res.Destination)
	for _, s := range res.Steps {
		fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", s)
	}
	digest, err := res.Image.Digest()
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "preview digest: %s\n", digest)
	return nil
}

var preview = previewCmd
