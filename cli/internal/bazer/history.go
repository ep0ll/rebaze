package bazer

import (
	"encoding/json"
	"fmt"

	hist "github.com/ep0ll/rebaze/internal/history"
	"github.com/spf13/cobra"
)

var historyFile string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show local mutation / copy / rebase history",
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVar(&historyFile, "path", "", "Override history file path")
}

func runHistory(cmd *cobra.Command, args []string) error {
	path := historyFile
	var err error
	if path == "" {
		path, err = hist.DefaultPath()
		if err != nil {
			return err
		}
	}
	ev, err := hist.Load(path)
	if err != nil {
		return err
	}
	if len(ev) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no history entries")
		return nil
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(ev)
}
