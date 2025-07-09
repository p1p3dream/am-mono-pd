package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const distsyncLockId = "ssm-params"

var mainCmd = &cobra.Command{
	Use:          "ssm-params",
	SilenceUsage: true,
}

func main() {
	mainCmd.PersistentFlags().String("lock-table", os.Getenv("ABODEMINE_LOCK_TABLE"), "The dynamodb lock table to use.")
	if err := viper.BindPFlag("lock-table", mainCmd.PersistentFlags().Lookup("lock-table")); err != nil {
		panic(err)
	}

	mainCmd.PersistentFlags().String("namespace", os.Getenv("ABODEMINE_NAMESPACE"), "The namespace to use for the SSM parameters.")
	if err := viper.BindPFlag("namespace", mainCmd.PersistentFlags().Lookup("namespace")); err != nil {
		panic(err)
	}

	if err := mainCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
