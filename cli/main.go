package main

import (
	"fmt"
	"os"

	"github.com/ep0ll/rebaze/cli/internal/bazer"
)

func main() {
	if err := bazer.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
