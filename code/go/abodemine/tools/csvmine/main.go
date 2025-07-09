package main

import (
	"os"

	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:          "csvmine",
	SilenceUsage: true,
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
