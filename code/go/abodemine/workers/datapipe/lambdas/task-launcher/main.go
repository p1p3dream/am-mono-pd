package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	"abodemine/lib/app"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/logging"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/domains/lambda"
)

func main() {
	aws_lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	defer logging.FlushSentry()
	if err := handlerFunc(ctx, event); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to handle SQS event.")
		return err
	}
	return nil
}

func handlerFunc(ctx context.Context, event events.SQSEvent) error {
	if err := gconf.LoadZerolog("info", true); err != nil {
		return errors.Forward(err, "3c0d1f2c-a1a7-460d-a3d1-0d108bd5f25e")
	}

	requestId, err := val.NewUUID7()
	if err != nil {
		return errors.Forward(err, "16da455f-8f21-42ab-a139-159b3828f439")
	}

	log.Info().
		Str("build_id", app.BuildId()).
		Str("build_version", app.BuildVersion()).
		Str("request_id", requestId.String()).
		Send()

	config, err := conf.ResolveAndLoad(ctx, "")
	if err != nil {
		return errors.Forward(err, "dbf57f7a-8caf-4c98-8ab9-f4a3e7f55a8a")
	}

	requestDomain := arc.NewDomain(&arc.NewDomainInput{})
	lambdaDomain := lambda.NewDomain(&lambda.NewDomainInput{
		Config: config,
	})

	r, err := requestDomain.CreateRequest(&arc.CreateRequestInput{
		Id:      requestId,
		Context: ctx,
	})
	if err != nil {
		return errors.Forward(err, "f132384b-693d-4f9b-aaab-9f56b1f71256")
	}

	in := &lambda.HandleTaskLauncherLambdaEventInput{
		SqsMessages: event.Records,
	}

	_, err = lambdaDomain.HandleTaskLauncherLambdaEvent(r, in)
	if err != nil {
		return errors.Forward(err, "aa70e3d9-6501-4c98-a1e5-031f22509ee3")
	}

	return nil
}
