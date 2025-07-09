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
	"abodemine/lib/storage"
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
			return errors.Forward(err, "27efc9f8-485d-4d20-b918-c98143706463")
		}

		requestId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "ca5bf013-8cd6-4494-a88e-df707b000ccc")
		}

		log.Info().
			Str("build_id", app.BuildId()).
			Str("build_version", app.BuildVersion()).
			Str("request_id", requestId.String()).
			Send()

		partnerId, err := uuid.Parse(viper.GetString("run.partner-id"))
		if err != nil {
			return &errors.Object{
				Id:     "52ab2941-0df0-40c7-868e-2b73cb722417",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Invalid PartnerId.",
			}
		}

		arcDomain := arc.NewDomain(&arc.NewDomainInput{
			DeploymentEnvironment: config.File.DeploymentEnvironment,
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
			return errors.Forward(err, "3d87fc94-058f-4474-ab5c-9613ccfac5fd")
		}

		var storageBackend storage.Backend

		switch {
		case viper.GetString("run.bucket") != "":
			storageBackend = &storage.S3Backend{
				AWS:    config.AWS.Get("default"),
				Bucket: viper.GetString("run.bucket"),
			}
		case viper.GetString("run.dir") != "":
			storageBackend = &storage.LocalBackend{
				FilesystemPath: viper.GetString("run.dir"),
			}
		default:
			return &errors.Object{
				Id:     "e94c0a56-674c-44f6-877c-92ca1cb6cbb6",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Missing run source.",
			}
		}

		workerId, err := val.NewUUID4()
		if err != nil {
			return errors.Forward(err, "42dccafd-5ba9-4106-8879-067b17647989")
		}

		log.Info().
			Str("partner_id", partnerId.String()).
			Str("bucket", viper.GetString("run.bucket")).
			Str("dir", viper.GetString("run.dir")).
			Bool("no_lock", viper.GetBool("run.no-lock")).
			Str("prefix", viper.GetString("run.prefix")).
			Int("file_buffer_size", config.File.FileBufferSize).
			Str("version", viper.GetString("run.version")).
			Str("worker_id", workerId.String()).
			Msg("Running loader.")

		_, err = workerDomain.ProcessDataSource(
			r,
			&worker.ProcessDataSourceInput{
				PartnerId:         partnerId,
				Backend:           storageBackend,
				FileBufferSize:    config.File.FileBufferSize,
				DatabaseBatchSize: 1000,
				PathPrefix:        viper.GetString("run.prefix"),
				PriorityGroup:     viper.GetInt32("run.priority-group"),
				NoLock:            viper.GetBool("run.no-lock"),
				Version:           viper.GetString("run.version"),
				WorkerId:          &workerId,
			},
		)
		if err != nil {
			return errors.Forward(err, "4608f372-102a-412d-9667-92c3b3dec44b")
		}

		return nil
	},
}

func init() {
	runCmd.PersistentFlags().String("bucket", "", "Bucket with data sources.")
	if err := viper.BindPFlag("run.bucket", runCmd.PersistentFlags().Lookup("bucket")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("dir", "", "Filesystem dir with data sources.")
	if err := viper.BindPFlag("run.dir", runCmd.PersistentFlags().Lookup("dir")); err != nil {
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

	runCmd.PersistentFlags().String("prefix", "", "Path prefix for the bucket.")
	if err := viper.BindPFlag("run.prefix", runCmd.PersistentFlags().Lookup("prefix")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().Int32("priority-group", 0, "Priority group to process.")
	if err := viper.BindPFlag("run.priority-group", runCmd.PersistentFlags().Lookup("priority-group")); err != nil {
		panic(err)
	}

	runCmd.PersistentFlags().String("version", "v2", "Loader version to use.")
	if err := viper.BindPFlag("run.version", runCmd.PersistentFlags().Lookup("version")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(runCmd)
}
