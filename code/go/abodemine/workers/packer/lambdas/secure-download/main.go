package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	"abodemine/lib/app"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/logging"
	"abodemine/lib/val"
	"abodemine/projects/packer/conf"
	"abodemine/projects/packer/domains/lambda"
)

func main() {
	aws_lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	defer logging.FlushSentry()

	resp, err := handlerFunc(ctx, event)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to handle APIGateway event.")
		return nil, err
	}

	return resp, nil
}

func handlerFunc(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	if err := gconf.LoadZerolog("info", true); err != nil {
		return nil, errors.Forward(err, "10e19239-1400-4661-9cdc-5adc0e26601d")
	}

	requestId, err := val.NewUUID7()
	if err != nil {
		return nil, errors.Forward(err, "7eecfc89-ba1a-4e2c-9077-051f5fd469ff")
	}

	log.Info().
		Str("build_id", app.BuildId()).
		Str("build_version", app.BuildVersion()).
		Str("request_id", requestId.String()).
		Send()

	config, err := conf.ResolveAndLoad(ctx, "")
	if err != nil {
		return nil, errors.Forward(err, "34e00bae-a2da-4bda-b004-efce74f664f4")
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
		return nil, errors.Forward(err, "0bd52847-8e39-4517-9ba0-d550e8367808")
	}

	in := &lambda.HandleSecureDownloadLambdaEventInput{
		Event: &event,
	}

	handleEventOut, err := lambdaDomain.HandleSecureDownloadLambdaEvent(r, in)
	if err != nil {
		return nil, errors.Forward(err, "654fb816-680a-4ce2-9a92-916db83ab48e")
	}

	out := &events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": handleEventOut.PresignedURL,
		},
	}

	return out, nil
}
