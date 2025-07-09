package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"abodemine/lib/distsync"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/val"
)

var getCmd = &cobra.Command{
	Use:          "get",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gconf.LoadZerolog("", false); err != nil {
			return errors.Forward(err, "f3dd57dd-fd93-4cd0-8be8-5a6ce3529791")
		}

		ctx := context.Background()

		awsConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return &errors.Object{
				Id:     "6b9b55f9-2757-440e-b310-c807df574968",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to load AWS configuration.",
				Cause:  err.Error(),
			}
		}

		ttl := time.Minute

		lock := &distsync.Lock{
			Id:   distsyncLockId,
			Type: distsync.LockTypeRead,
			Ttl:  ttl,
		}

		locker, err := getLock(&getLockInput{
			AWS:       awsConfig,
			Ctx:       ctx,
			Lock:      lock,
			LockTable: viper.GetString("lock-table"),
			Timeout:   time.Minute,
		})
		if err != nil {
			return errors.Forward(err, "e783a0cd-f143-4663-a07d-e6c83e8aefaa")
		}

		defer func() {
			if err := locker.Unlock(ctx); err != nil {
				log.Error().
					Str("id", "94cdd9c9-8b6f-4180-b84a-dd8d89cd4de7").
					Err(err).
					Send()
			}
		}()

		if _, err := getParams(&getParamsInput{
			AWS:         awsConfig,
			Context:     ctx,
			Format:      viper.GetString("get.format"),
			KeepSuffix:  viper.GetBool("get.keep-suffix"),
			Namespace:   viper.GetString("namespace"),
			Out:         viper.GetString("get.out"),
			Prefix:      viper.GetString("get.prefix"),
			StripPrefix: viper.GetString("get.strip-prefix"),
		}); err != nil {
			return errors.Forward(err, "8a981a15-038f-4956-b266-10a3e1dac5c3")
		}

		return nil
	},
}

func init() {
	// When saving the parameters locally to bulk edit them and later push
	// them back to SSM, we must keep the suffixes in the parameter names.
	getCmd.PersistentFlags().Bool("keep-suffix", false, "Do not strip suffixes from parameter names.")
	if err := viper.BindPFlag("get.keep-suffix", getCmd.PersistentFlags().Lookup("keep-suffix")); err != nil {
		panic(err)
	}

	getCmd.PersistentFlags().String("format", "json", "Output format: json, yaml.")
	if err := viper.BindPFlag("get.format", getCmd.PersistentFlags().Lookup("format")); err != nil {
		panic(err)
	}

	getCmd.PersistentFlags().String("out", "", "File path to write the parameters to.")
	if err := viper.BindPFlag("get.out", getCmd.PersistentFlags().Lookup("out")); err != nil {
		panic(err)
	}

	getCmd.PersistentFlags().String("prefix", "", "Filter parameters by prefix.")
	if err := viper.BindPFlag("get.prefix", getCmd.PersistentFlags().Lookup("prefix")); err != nil {
		panic(err)
	}

	getCmd.PersistentFlags().String("strip-prefix", "", "Filter parameters by prefix.")
	if err := viper.BindPFlag("get.strip-prefix", getCmd.PersistentFlags().Lookup("strip-prefix")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(getCmd)
}

type getParamsInput struct {
	AWS          aws.Config
	Context      context.Context
	Format       string
	KeepSuffix   bool
	Namespace    string
	Out          string
	Prefix       string
	ReturnParams bool
	StripPrefix  string
}

type getParamsOutput struct {
	Params map[string]*ssmParam
}

func getParams(in *getParamsInput) (*getParamsOutput, error) {
	namespacePrefix := "/" + in.Namespace + "/"

	if in.Prefix == "" {
		in.Prefix = namespacePrefix
	}

	if in.StripPrefix == "" && in.Prefix == namespacePrefix {
		in.StripPrefix = namespacePrefix
	}

	nextToken := val.PtrRef("")
	params := make(map[string]*ssmParam)
	ssmClient := ssm.NewFromConfig(in.AWS)

	for {
		getParametersOut, err := ssmClient.GetParametersByPath(in.Context, &ssm.GetParametersByPathInput{
			Path:      val.PtrRef(in.Prefix),
			Recursive: val.PtrRef(true),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, &errors.Object{
				Id:     "28d327d8-0286-495f-bbd7-45895844667b",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to get parameters.",
				Cause:  err.Error(),
			}
		}

		for _, param := range getParametersOut.Parameters {
			if in.StripPrefix != "" {
				param.Name = val.PtrRef(strings.TrimPrefix(*param.Name, in.StripPrefix))
			}

			paramName := *param.Name
			ssmParam := &ssmParam{
				HaveName: paramName,
				Type:     param.Type,
			}

			if len(paramName) > 5 {
				idx := len(paramName) - 5
				suffix := paramName[idx:]

				if suffix == ".JSON" {
					ssmParam.Encoding = ParamEncodingOptionJSON
					ssmParam.WantName = paramName[:idx]

					var v any

					if err := json.Unmarshal([]byte(*param.Value), &v); err != nil {
						return nil, &errors.Object{
							Id:     "d58ef19c-4d6d-4f23-854f-1703b32d6e22",
							Code:   errors.Code_UNKNOWN,
							Detail: "Failed to unmarshal JSON value.",
							Cause:  err.Error(),
							Meta: map[string]any{
								"value": *param.Value,
							},
						}
					}

					ssmParam.Value = v
				}
			}

			if ssmParam.Encoding == ParamEncodingOptionNone {
				ssmParam.WantName = paramName
				ssmParam.Value = *param.Value
			}

			params[paramName] = ssmParam
		}

		if getParametersOut.NextToken == nil {
			break
		}

		nextToken = getParametersOut.NextToken
	}

	out := &getParamsOutput{}

	if in.ReturnParams {
		out.Params = params
		return out, nil
	}

	var err error
	var marshaledOutput []byte

	switch in.Format {
	case "json":
		marshaledOutput, err = json.Marshal(params)
	case "yaml":
		buf := new(bytes.Buffer)

		encoder := yaml.NewEncoder(buf)
		encoder.SetIndent(2)

		err = encoder.Encode(params)
		marshaledOutput = buf.Bytes()
	default:
		return nil, &errors.Object{
			Id:     "09f6a85a-32ed-4dd7-8cc3-f9b4ed330ce6",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported format.",
			Meta: map[string]any{
				"format": in.Format,
			},
		}
	}

	if err != nil {
		return nil, &errors.Object{
			Id:     "7520f0a1-7f4a-4aba-b314-936008e9e2ca",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to marshal parameters.",
			Cause:  err.Error(),
			Meta: map[string]any{
				"format": in.Format,
			},
		}
	}

	if in.Out == "" || in.Out == "-" {
		fmt.Println(string(marshaledOutput))
		return out, nil
	}

	if err := os.MkdirAll(filepath.Dir(in.Out), 0755); err != nil {
		return nil, &errors.Object{
			Id:     "f9fc8e3a-3291-4e3c-94e6-d3ec4de191c3",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create output directory.",
			Cause:  err.Error(),
		}
	}

	if err := os.WriteFile(in.Out, marshaledOutput, 0644); err != nil {
		return nil, &errors.Object{
			Id:     "45690cc3-1057-4001-a216-e32b342bf113",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to write output.",
			Cause:  err.Error(),
		}
	}

	return out, nil
}
