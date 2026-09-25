package bazer

import (
	"encoding/json"
	"fmt"

	"github.com/ep0ll/rebaze/internal/history"
	"github.com/spf13/cobra"
)

var historyPath string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show local mutation / copy / rebase history",
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVar(&historyPath, "path", "", "Override history file path")
}

func runHistory(cmd *cobra.Command, args []string) error {
	path := historyPath
	var err error
	if path == "" {
		path, err = history.DefaultPath()
		if err != nil {
			return err
		}
	}
	ev, err := history.Load(path)
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

var history = historyCmd
