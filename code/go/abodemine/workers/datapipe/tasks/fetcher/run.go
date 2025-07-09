package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"abodemine/domains/arc"
	"abodemine/lib/app"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/domains/worker"
)

var runCmd = &cobra.Command{
	Use:          "run",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		config, err := conf.ResolveAndLoad(ctx, viper.GetString("config"))
		if err != nil {
			return errors.Forward(err, "c3597379-23e6-4185-bc18-07eecb5e4734")
		}

		requestId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "0ba6ca58-a8c7-4a3d-ad00-3fad927251a7")
		}

		log.Info().
			Str("build_id", app.BuildId()).
			Str("build_version", app.BuildVersion()).
			Str("request_id", requestId.String()).
			Send()

		partnerId, err := uuid.Parse(viper.GetString("run.partner-id"))
		if err != nil {
			return &errors.Object{
				Id:     "e2262439-7c81-4655-bd45-5b368a0f6069",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid PartnerId.",
			}
		}

		arcDomain := arc.NewDomain(&arc.NewDomainInput{
			DeploymentEnvironment: config.File.DeploymentEnvironment,
			PgxPool:               config.PgxPool,
		})

		workerDomain := worker.NewDomain(&worker.NewDomainInput{
			Config: config,
		})

		r, err := arcDomain.CreateRequest(&arc.CreateRequestInput{
			Id:      requestId,
			Context: ctx,
		})
		if err != nil {
			return errors.Forward(err, "cbb141ef-f8a4-4b41-bc82-536ba6ba0fa7")
		}

		log.Info().
			Str("partner_id", partnerId.String()).
			Msg("Running fetcher.")

		_, err = workerDomain.FetchDataSource(
			r,
			&worker.FetchDataSourceInput{
				PartnerId:         partnerId,
				RcloneCheckers:    viper.GetString("run.rclone-checkers"),
				RcloneTransfers:   viper.GetString("run.rclone-transfers"),
				RcloneDestination: viper.GetString("run.rclone-dst"),
				RcloneSource:      viper.GetString("run.rclone-src"),
			},
		)
		if err != nil {
			return errors.Forward(err, "813bbc7a-d573-4bc5-9892-6036b9be7926")
		}

		return nil
	},
}

func init() {
	runCmd.PersistentFlags().String("partner-id", "", "Id of the data partner.")
	if err := viper.BindPFlag("run.partner-id", runCmd.PersistentFlags().Lookup("partner-id")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("rclone-dst", "", "Rclone destination, with path.")
	if err := viper.BindPFlag("run.rclone-dst", runCmd.PersistentFlags().Lookup("rclone-dst")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("rclone-src", "", "Rclone source, with path.")
	if err := viper.BindPFlag("run.rclone-src", runCmd.PersistentFlags().Lookup("rclone-src")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("rclone-checkers", "8", "Rclone checkers count.")
	if err := viper.BindPFlag("run.rclone-checkers", runCmd.PersistentFlags().Lookup("rclone-checkers")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("rclone-transfers", "8", "Rclone transfers count.")
	if err := viper.BindPFlag("run.rclone-transfers", runCmd.PersistentFlags().Lookup("rclone-transfers")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(runCmd)
}
