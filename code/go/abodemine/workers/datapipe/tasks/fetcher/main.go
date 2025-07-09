package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"abodemine/lib/logging"
)

var mainCmd = &cobra.Command{
	Use:          "worker",
	SilenceUsage: true,
}

func init() {
	mainCmd.PersistentFlags().String("config", "", "Path to config file.")
	if err := viper.BindPFlag("config", mainCmd.PersistentFlags().Lookup("config")); err != nil {
		panic(err)
	}
}

func main() {
	logging.ExecuteCobraCommand(mainCmd)
}
