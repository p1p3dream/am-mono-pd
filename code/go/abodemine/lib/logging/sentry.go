package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryzerolog "github.com/getsentry/sentry-go/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"abodemine/lib/app"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/val"
)

// zerologSentryWriter is a zerolog writer that
// processes logs before sending to Sentry.
type zerologSentryWriter struct {
	w zerolog.LevelWriter
}

func (w *zerologSentryWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *zerologSentryWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	return w.w.WriteLevel(level, p)
}

// InitSentry initializes the Sentry client for error tracking
// and configures it to capture logs and panics.
func InitSentry(conf *gconf.Sentry, deploymentEnvironment int, noLogColor bool) (*sentryzerolog.Writer, error) {
	if conf.DSN == "" {
		return nil, &errors.Object{
			Id:     "723ee5c3-f835-4443-ae8b-e08a519d18f1",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Sentry DSN is not set.",
		}
	}

	// TODO: Find if there's a better way than marshal/unmarshal to get the chain.
	beforeSend := func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		if len(event.Exception) == 0 {
			return event
		}

		exception := &event.Exception[0]
		chainStr, err := strconv.Unquote(`"` + exception.Value + `"`)
		if err != nil {
			fmt.Printf("ERROR: Failed to unquote exception value: %s.\n", err)
			return event
		}

		chain := &errors.Chain{}

		if err := json.Unmarshal([]byte(chainStr), chain); err != nil {
			fmt.Printf("ERROR: Failed to unmarshal chain: %s.\n", err)
			return event
		}

		first := chain.First()
		exception.Value = val.Coalesce(
			first.Detail,
			first.Label,
			first.Cause,
			first.Id,
		)
		event.Extra["chain"] = chainStr

		return event
	}

	deploymentEnvironmentStr := strings.ToLower(gconf.DeploymentEnvironmentToString(deploymentEnvironment))

	// Initialize Sentry with the provided DSN.
	err := sentry.Init(sentry.ClientOptions{
		Dsn: conf.DSN,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		SendDefaultPII:   true,
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		AttachStacktrace: true,
		Environment:      deploymentEnvironmentStr,
		Dist:             app.BuildId(),
		Release:          app.BuildVersion(),
		BeforeSend:       beforeSend,
	})
	if err != nil {
		return nil, &errors.Object{
			Id:     "9e1788c4-d6a4-42cd-9a55-6479602fc7c0",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Failed to initialize Sentry.",
			Cause:  err.Error(),
		}
	}

	sentryWriter, err := sentryzerolog.New(sentryzerolog.Config{
		ClientOptions: sentry.ClientOptions{
			Dsn:              conf.DSN,
			SendDefaultPII:   true,
			AttachStacktrace: true,
			Environment:      deploymentEnvironmentStr,
			Dist:             app.BuildId(),
			Release:          app.BuildVersion(),
			BeforeSend:       beforeSend,
		},
		Options: sentryzerolog.Options{
			Levels: []zerolog.Level{
				zerolog.WarnLevel,
				zerolog.ErrorLevel,
				zerolog.FatalLevel,
				zerolog.PanicLevel,
			},
			FlushTimeout:    3 * time.Second,
			WithBreadcrumbs: true,
		},
	})
	if err != nil {
		return nil, &errors.Object{
			Id:     "11419245-1409-4861-9852-99ac4402f210",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Failed to create Sentry writer.",
			Cause:  err.Error(),
		}
	}
	// defer sentryWriter.Close()

	// Configure zerolog to include caller information.
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", file, line)
	}
	zerolog.CallerSkipFrameCount = 2 // Skip frames to get to the actual caller.

	// Configure how errors are marshaled when using Err().
	zerolog.ErrorMarshalFunc = func(err error) any {
		chain := errors.AsChain(err)
		b, err := json.Marshal(chain)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to marshal error chain.")
			return ""
		}
		return string(b)
	}

	zerologConfig, err := gconf.GetZerologConfig("", noLogColor)
	if err != nil {
		return nil, errors.Forward(err, "a04a834a-6bdf-42f1-89ef-b06abab81446")
	}

	var consoleWriter io.Writer

	if zerologConfig.ConsoleWriter.NoColor {
		consoleWriter = zerologConfig.ConsoleWriter.Out
	} else {
		consoleWriter = zerologConfig.ConsoleWriter
	}

	customWriter := &zerologSentryWriter{w: sentryWriter}

	// Configure the global logger to use both console and sentry writers.
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(
		consoleWriter,
		customWriter,
	)).
		With().
		Caller().
		Timestamp().
		Logger()

	log.Info().
		Str("environment", deploymentEnvironmentStr).
		Msg("Sentry initialized successfully")

	return sentryWriter, nil
}

// FlushSentry ensures all Sentry events are sent before the program exits.
func FlushSentry() {
	sentry.Flush(3 * time.Second)
}

// CaptureException sends an error to Sentry
func CaptureException(err error) {
	sentry.CaptureException(err)
}

// CaptureMessage sends a message to Sentry
func CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}

func CaptureEvent(event *sentry.Event) {
	sentry.CaptureEvent(event)
}

// ExecuteCobraCommand executes a Cobra command and
// ensures Sentry is flushed on exit.
func ExecuteCobraCommand(cmd *cobra.Command) {
	defer FlushSentry()

	if err := cmd.Execute(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to execute command.")

		// Flushing manually because on os.Exit,
		// deferred functions are not executed.
		FlushSentry()

		os.Exit(-1)
	}
}
