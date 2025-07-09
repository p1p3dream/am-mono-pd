package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/lib/app"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/domains/worker"
	"abodemine/repositories/opensearch"
)

var runCmd = &cobra.Command{
	Use:          "run",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		config, err := conf.ResolveAndLoad(ctx, viper.GetString("config"))
		if err != nil {
			return errors.Forward(err, "85b9a885-3bb1-49d4-8ea5-075fd7604c36")
		}

		requestId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "03171f81-d877-4711-9c76-edf960aadfe9")
		}

		log.Info().
			Str("build_id", app.BuildId()).
			Str("build_version", app.BuildVersion()).
			Str("request_id", requestId.String()).
			Send()

		arcDomain := arc.NewDomain(&arc.NewDomainInput{
			DeploymentEnvironment: config.File.DeploymentEnvironment,
			OpenSearch:            config.OpenSearch,
			PgxPool:               config.PgxPool,
		})

		addressDomain := address.NewDomain(&address.NewDomainInput{})

		workerDomain := worker.NewDomain(&worker.NewDomainInput{
			Config:        config,
			AddressDomain: addressDomain,
			OsSearchRepository: opensearch.NewRepository(&opensearch.NewRepositoryInput{
				ConfigKey: consts.ConfigKeyOpenSearchSearch,
			}),
		})

		r, err := arcDomain.CreateRequest(&arc.CreateRequestInput{
			Id:      requestId,
			Context: ctx,
		})
		if err != nil {
			return errors.Forward(err, "4bedae33-14dd-4215-a2db-33977ab2be4b")
		}

		workerId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "6506c7c6-9ba2-470e-8d5e-bb4a8fb54106")
		}

		log.Info().
			Bool("no_lock", viper.GetBool("run.no-lock")).
			Str("version", viper.GetString("run.version")).
			Str("worker_id", workerId.String()).
			Msg("Running synther.")

		_, err = workerDomain.SyncProperties(
			r,
			&worker.SyncPropertiesInput{
				NoLock:              viper.GetBool("run.no-lock"),
				OpenSearchIndexName: viper.GetString("run.os-index-name"),
				Version:             viper.GetString("run.version"),
				WorkerId:            &workerId,
			},
		)
		if err != nil {
			return errors.Forward(err, "75184716-a6a7-41a7-a516-fe2eccad67e9")
		}

		return nil
	},
}

func init() {
	runCmd.PersistentFlags().String("os-index-name", "", "Name of the opensearch index to use.")
	if err := viper.BindPFlag("run.os-index-name", runCmd.PersistentFlags().Lookup("os-index-name")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().Bool("no-lock", false, "Do not check for locks.")
	if err := viper.BindPFlag("run.no-lock", runCmd.PersistentFlags().Lookup("no-lock")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("version", "v1", "Synther version to use.")
	if err := viper.BindPFlag("run.version", runCmd.PersistentFlags().Lookup("version")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(runCmd)
}
