package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"abodemine/lib/app"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/projects/saas/conf"
	"abodemine/projects/saas/handlers"
)

var listenCmd = &cobra.Command{
	Use:          "listen",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().
			Str("build_id", app.BuildId()).
			Str("build_version", app.BuildVersion()).
			Send()

		config, err := conf.ResolveAndLoad(viper.GetString("config"))
		if err != nil {
			return errors.Forward(err, "f4ab0add-2f26-41df-aa30-47dee60ae4c6")
		}

		log.Info().
			Int("pid", os.Getpid()).
			Msg("Listening")

		wg := new(sync.WaitGroup)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := startHTTP(config); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to start HTTP server")

				panic(err)
			}
		}()

		wg.Wait()

		return nil
	},
}

func init() {
	mainCmd.AddCommand(listenCmd)
}

func startHTTP(config *conf.Config) error {
	httpServer := config.File.HttpServer

	if httpServer == nil {
		return errors.New("missing config")
	}

	if httpServer.Tls == nil {
		return errors.New("missing TLS")
	}

	tc, err := gconf.GetTLS(httpServer.Tls)
	if err != nil {
		return fmt.Errorf("failed to get tls: %w", err)
	}

	bind := fmt.Sprintf(
		"%s:%d",
		httpServer.Bind,
		httpServer.Port,
	)

	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("failed to net listen: %w", err)
	}

	router := handlers.Router(config)

	server := &http.Server{
		Handler:   router,
		TLSConfig: tc,
	}

	log.Info().
		Str("bind", bind).
		Strs("hostnames", httpServer.Hostnames).
		Msg("Serving HTTPS.")

	return server.ServeTLS(listener, "", "")
}
