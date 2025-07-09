package main

import (
	"context"

	"github.com/google/uuid"
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
			return errors.Forward(err, "a0d2d114-f337-4d08-b386-722f4c91bfc8")
		}

		requestId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "a8fdbc7b-131f-4d13-a2d5-9744b3547b74")
		}

		log.Info().
			Str("build_id", app.BuildId()).
			Str("build_version", app.BuildVersion()).
			Str("request_id", requestId.String()).
			Send()

		arcDomain := arc.NewDomain(&arc.NewDomainInput{
			DeploymentEnvironment: config.File.DeploymentEnvironment,
			Flags:                 config.File.Flags,
			OpenSearch:            config.OpenSearch,
			PgxPool:               config.PgxPool,
		})

		partnerId, err := uuid.Parse(viper.GetString("run.partner-id"))
		if err != nil {
			return &errors.Object{
				Id:     "d3ce6c56-3885-4b8e-9f79-ebf9def3e131",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid PartnerId.",
			}
		}

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
			return errors.Forward(err, "98580d52-82aa-4413-8cc3-68e5989f9a36")
		}

		log.Info().
			Int("file_buffer_size", config.File.FileBufferSize).
			Str("index_name", viper.GetString("run.index-name")).
			Bool("no_lock", viper.GetBool("run.no-lock")).
			Str("partner_id", partnerId.String()).
			Str("version", viper.GetString("run.version")).
			Msg("Running OpenSearch loader.")

		_, err = workerDomain.LoadOpenSearch(r, &worker.LoadOpenSearchInput{
			BatchSize:      1000,
			FileBufferSize: config.File.FileBufferSize,
			IndexName:      viper.GetString("run.index-name"),
			NoLock:         viper.GetBool("run.no-lock"),
			PartnerId:      partnerId,
			Version:        viper.GetString("run.version"),
		})
		if err != nil {
			return errors.Forward(err, "387fbca9-3b25-45ee-9316-00721fd3db7b")
		}

		return nil
	},
}

func init() {
	runCmd.PersistentFlags().String("index-name", "", "Name of the opensearch index to use.")
	if err := viper.BindPFlag("run.index-name", runCmd.PersistentFlags().Lookup("index-name")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().Bool("no-lock", false, "Do not check for locks.")
	if err := viper.BindPFlag("run.no-lock", runCmd.PersistentFlags().Lookup("no-lock")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("partner-id", "", "Id of the data partner.")
	if err := viper.BindPFlag("run.partner-id", runCmd.PersistentFlags().Lookup("partner-id")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("version", "v1", "Loader version to use.")
	if err := viper.BindPFlag("run.version", runCmd.PersistentFlags().Lookup("version")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(runCmd)
}
