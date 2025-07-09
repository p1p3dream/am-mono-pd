package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/rs/zerolog/log"

	"abodemine/lib/distsync"
	"abodemine/lib/errors"
)

type getLockInput struct {
	AWS       aws.Config
	Ctx       context.Context
	Lock      *distsync.Lock
	LockTable string
	Timeout   time.Duration
}

func getLock(in *getLockInput) (distsync.Locker, error) {
	log.Info().Msg("Acquiring lock.")

	locker := &distsync.DynamoDB{
		PollInterval: time.Second,
		Client:       dynamodb.NewFromConfig(in.AWS),
		TableName:    in.LockTable,
	}

	ctx, cancel := context.WithTimeout(in.Ctx, in.Timeout)
	defer cancel()

	if err := locker.Lock(ctx, in.Lock); err != nil {
		return nil, errors.Forward(err, "55e0781e-1d18-4c88-9287-af16cb709b30")
	}

	log.Info().Msg("Lock acquired.")

	return locker, nil
}

type ParamEncodingOption int

const (
	ParamEncodingOptionNone ParamEncodingOption = 0
	ParamEncodingOptionJSON ParamEncodingOption = 1
)

type ssmParam struct {
	// Wether the param has an encoded value or not, and what type it is.
	Encoding ParamEncodingOption

	// HaveName is the (actual) name the param has on SSM, i.e.,
	// the name you would use to get the param.
	HaveName string

	// WantName is the name the param MUST have in the output, for
	// usage during config evaluation.
	// HaveName and WantName will probably differ when
	// ParamEncodingOption is not ParamEncodingOptionNone.
	WantName string

	// The param type on SSM.
	Type types.ParameterType

	Value any

	// Used to indicate the param MUST NOT be deleted
	// during config updates.
	doNotDelete bool
}

func (p *ssmParam) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *ssmParam) MarshalYAML() (any, error) {
	return p.Value, nil
}

type fileToSsmParamsInput struct {
	ParamsFile map[string]any
}

type fileToSsmParamsOutput struct {
	Params map[string]*ssmParam
}

func fileToSsmParams(in *fileToSsmParamsInput) (*fileToSsmParamsOutput, error) {
	params := make(map[string]*ssmParam)

	for k, v := range in.ParamsFile {
		ssmParam := &ssmParam{
			HaveName: k,
			Value:    v,
			Type:     types.ParameterTypeString,
		}

		if len(k) > 5 {
			idx := len(k) - 5
			suffix := k[idx:]

			if suffix == ".JSON" {
				ssmParam.Encoding = ParamEncodingOptionJSON
				ssmParam.WantName = k[:idx]
			}
		}

		if ssmParam.Encoding == ParamEncodingOptionNone {
			ssmParam.WantName = k
		}

		params[k] = ssmParam
	}

	out := &fileToSsmParamsOutput{
		Params: params,
	}

	return out, nil
}
