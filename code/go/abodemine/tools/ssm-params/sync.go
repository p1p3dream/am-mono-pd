package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"abodemine/lib/distsync"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/val"
)

var syncCmd = &cobra.Command{
	Use:          "sync",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gconf.LoadZerolog("", false); err != nil {
			return errors.Forward(err, "982241ac-aef1-4eff-997d-4cc03a29a7be")
		}

		ctx := context.Background()

		awsConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return &errors.Object{
				Id:     "91935818-dfb0-4680-a326-fa78f5a54c4f",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to load AWS configuration.",
				Cause:  err.Error(),
			}
		}

		ttl := time.Minute

		lock := &distsync.Lock{
			Id:   distsyncLockId,
			Type: distsync.LockTypeWrite,
			Ttl:  ttl,
		}

		locker, err := getLock(&getLockInput{
			AWS:       awsConfig,
			Ctx:       ctx,
			Lock:      lock,
			LockTable: viper.GetString("lock-table"),
			Timeout:   ttl - 5*time.Second,
		})
		if err != nil {
			return errors.Forward(err, "552f7a84-f6fa-40bd-aaa4-ee458dc0cc3d")
		}

		defer func() {
			if err := locker.Unlock(ctx); err != nil {
				log.Error().
					Str("id", "e654778a-8947-4a52-82e0-ebda78074f75").
					Err(err).
					Send()
			}
		}()

		if err := syncParams(&syncParamsInput{
			AWS:         awsConfig,
			Ctx:         ctx,
			LockTable:   viper.GetString("lock-table"),
			Namespace:   viper.GetString("namespace"),
			Prefix:      viper.GetString("sync.prefix"),
			UpdatesFile: viper.GetString("sync.updates-file"),
		}); err != nil {
			return errors.Forward(err, "e25adcdb-936a-4fd7-9eb5-ffb49ba82bf2")
		}

		return nil
	},
}

func init() {
	syncCmd.PersistentFlags().String("prefix", "", "Filter parameters by prefix.")
	if err := viper.BindPFlag("sync.prefix", syncCmd.PersistentFlags().Lookup("prefix")); err != nil {
		panic(err)
	}

	syncCmd.PersistentFlags().String("updates-file", "", "Path to the updates file.")
	if err := viper.BindPFlag("sync.updates-file", syncCmd.PersistentFlags().Lookup("updates-file")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(syncCmd)
}

type syncParamsInput struct {
	AWS         aws.Config
	Ctx         context.Context
	LockTable   string
	Namespace   string
	Prefix      string
	UpdatesFile string
}

func syncParams(in *syncParamsInput) error {
	if in.UpdatesFile == "" {
		return &errors.Object{
			Id:     "a678b920-31ed-43cc-960a-ee43c4e99a7c",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Updates file is required.",
		}
	}

	b, err := os.ReadFile(in.UpdatesFile)
	if err != nil {
		return &errors.Object{
			Id:     "8a74bc6e-3eb1-462d-a0d3-1ad53c985670",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to open updates file.",
			Cause:  err.Error(),
		}
	}

	newParamsFile := make(map[string]any)
	updatesFileExt := strings.ToLower(strings.TrimSpace(filepath.Ext(in.UpdatesFile)))

	switch updatesFileExt {
	case ".json":
		if err := json.Unmarshal(b, &newParamsFile); err != nil {
			return &errors.Object{
				Id:     "a3354f84-1ff1-4ea8-8fa5-1c2b13a91600",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to unmarshal updates file.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"file_ext": updatesFileExt,
				},
			}
		}
	case ".yaml":
		if err := yaml.Unmarshal(b, &newParamsFile); err != nil {
			return &errors.Object{
				Id:     "6500d892-1ce6-430e-9a24-166eb8c6bd56",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to unmarshal updates file.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"file_ext": updatesFileExt,
				},
			}
		}
	default:
		return &errors.Object{
			Id:     "0f161658-deb1-425b-9d42-faf7c77c5a96",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported file format",
			Meta: map[string]any{
				"file_ext": updatesFileExt,
			},
		}
	}

	fileToSsmParamsOut, err := fileToSsmParams(&fileToSsmParamsInput{
		ParamsFile: newParamsFile,
	})
	if err != nil {
		return errors.Forward(err, "5fe72f21-28f7-4226-bde5-62e36e707c8b")
	}

	newParams := fileToSsmParamsOut.Params

	log.Info().Msg("Getting current SSM parameters.")

	getParamsOut, err := getParams(&getParamsInput{
		AWS:          in.AWS,
		Context:      in.Ctx,
		Format:       "json",
		Namespace:    in.Namespace,
		Prefix:       in.Prefix,
		ReturnParams: true,
	})

	if err != nil {
		return errors.Forward(err, "0cb98bd5-a6ee-4943-8821-f6a721745997")
	}

	oldParams := getParamsOut.Params

	compareParamsOut, err := compareParams(&compareParamsInput{
		NewParams: newParams,
		OldParams: oldParams,
	})
	if err != nil {
		return errors.Forward(err, "64e06d11-c48d-49e9-a761-c696ed3109d1")
	}

	if len(compareParamsOut.UpdateCommands) == 0 {
		fmt.Println("No changes to apply.")
		return nil
	}

	fmt.Println("List of changes:")
	indent := "    "

	for _, cmd := range compareParamsOut.UpdateCommands {
		switch cmd.UpdateOption {
		case UpdateOptionCreate:
			fmt.Printf("%s%s %s\n", indent, color.GreenString("CREATE"), cmd.Param.HaveName)
		case UpdateOptionUpdate:
			fmt.Printf("%s%s %s\n", indent, color.YellowString("UPDATE"), cmd.Param.HaveName)
		case UpdateOptionDelete:
			fmt.Printf("%s%s %s\n", indent, color.RedString("DELETE"), cmd.Param.HaveName)
		}

		switch cmd.Param.Encoding {
		case ParamEncodingOptionNone:
			fmt.Printf("%sValue: %q\n", indent, cmd.Param.Value)
		case ParamEncodingOptionJSON:
			newValue, err := json.MarshalIndent(cmd.Param.Value, indent, indent)
			if err != nil {
				return &errors.Object{
					Id:     "670613dd-9c07-4f86-a6bc-8d176044cded",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to marshal JSON.",
					Cause:  err.Error(),
				}
			}

			if cmd.UpdateOption != UpdateOptionUpdate {
				fmt.Printf("%sValue:\n%s%s\n", indent, indent, string(newValue))
			} else {
				oldValue, err := json.MarshalIndent(cmd.Meta["old_value"], indent, indent)
				if err != nil {
					return &errors.Object{
						Id:     "1b4cf7fe-3ab1-4a5a-8f13-fcf55c94439b",
						Code:   errors.Code_UNKNOWN,
						Detail: "Failed to marshal JSON.",
						Cause:  err.Error(),
					}
				}

				diffFrom := string(oldValue)
				diffTo := string(newValue)
				dmp := diffmatchpatch.New()
				dmp.MatchThreshold = 0
				diffs := dmp.DiffMain(diffFrom, diffTo, true)
				diffs = dmp.DiffCleanupSemantic(diffs)

				fmt.Printf("%sDiff:\n%s%s\n", indent, indent, dmp.DiffPrettyText(diffs))
			}
		}

		fmt.Println()
	}

	fmt.Printf("\nDo you want to continue? (yes/no): ")

	var response string

	if _, err := fmt.Scanln(&response); err != nil {
		return &errors.Object{
			Id:     "736dbdd9-c4a4-4cb4-b263-12a5438eb17c",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to read user input.",
			Cause:  err.Error(),
		}
	}

	if strings.ToLower(response) != "yes" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	log.Info().Msg("Applying changes.")

	prefix := path.Join("/", in.Namespace)
	ssmClient := ssm.NewFromConfig(in.AWS)

	for _, cmd := range compareParamsOut.UpdateCommands {
		if err := applyUpdateCommand(&applyUpdateCommandInput{
			Command:   cmd,
			Context:   in.Ctx,
			Prefix:    prefix,
			SSMClient: ssmClient,
		}); err != nil {
			return errors.Forward(err, "5f760167-e09d-4506-b5ef-6f420ee57f6c")
		}
	}

	return nil
}

type UpdateOption int

const (
	UpdateOptionCreate UpdateOption = 0
	UpdateOptionUpdate UpdateOption = 1
	UpdateOptionDelete UpdateOption = 2
)

type updateCommand struct {
	UpdateOption UpdateOption
	Param        *ssmParam
	Description  string
	Meta         map[string]any
}

type compareParamsInput struct {
	NewParams map[string]*ssmParam
	OldParams map[string]*ssmParam
}

type compareParamsOutput struct {
	UpdateCommands []*updateCommand
}

func compareParams(in *compareParamsInput) (*compareParamsOutput, error) {
	var commands []*updateCommand

	for k, v := range in.NewParams {
		oldParam, exists := in.OldParams[k]

		if exists {
			oldParam.doNotDelete = true
		}

		switch {
		case !exists:
			commands = append(commands, &updateCommand{
				UpdateOption: UpdateOptionCreate,
				Param:        v,
			})
		case oldParam.Encoding != v.Encoding:
			commands = append(commands, &updateCommand{
				UpdateOption: UpdateOptionUpdate,
				Param:        v,
				Description:  "Encoding changed.",
				Meta: map[string]any{
					"new_encoding": v.Encoding,
					"old_encoding": oldParam.Encoding,
				},
			})
		case oldParam.Type != v.Type:
			commands = append(commands, &updateCommand{
				UpdateOption: UpdateOptionUpdate,
				Param:        v,
				Description:  "Type changed.",
				Meta: map[string]any{
					"new_type": v.Type,
					"old_type": oldParam.Type,
				},
			})
		case !reflect.DeepEqual(oldParam.Value, v.Value):
			commands = append(commands, &updateCommand{
				UpdateOption: UpdateOptionUpdate,
				Param:        v,
				Description:  "Value changed.",
				Meta: map[string]any{
					"new_value": v.Value,
					"old_value": oldParam.Value,
				},
			})
		}
	}

	for _, v := range in.OldParams {
		if v.doNotDelete {
			continue
		}

		commands = append(commands, &updateCommand{
			UpdateOption: UpdateOptionDelete,
			Param:        v,
		})
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Param.HaveName < commands[j].Param.HaveName
	})

	out := &compareParamsOutput{
		UpdateCommands: commands,
	}

	return out, nil
}

type applyUpdateCommandInput struct {
	Command   *updateCommand
	Context   context.Context
	Prefix    string
	SSMClient *ssm.Client
}

func applyUpdateCommand(in *applyUpdateCommandInput) error {
	param := in.Command.Param
	paramName := val.PtrRef(path.Join(in.Prefix, param.HaveName))
	ssmClient := in.SSMClient

	log.Info().
		Str("param_name", *paramName).
		Str("update_option", fmt.Sprintf("%v", in.Command.UpdateOption)).
		Msg("Updating SSM param.")

	var paramValue *string

	if in.Command.UpdateOption != UpdateOptionDelete {
		switch param.Encoding {
		case ParamEncodingOptionNone:
			paramValue = val.PtrRef(param.Value.(string))
		case ParamEncodingOptionJSON:
			b, err := json.Marshal(param.Value)
			if err != nil {
				return &errors.Object{
					Id:     "da7d0e82-c5bd-46cf-b467-8a7faa3e572a",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to marshal JSON.",
					Cause:  err.Error(),
				}
			}

			paramValue = val.PtrRef(string(b))
		}
	}

	switch in.Command.UpdateOption {
	case UpdateOptionCreate:
		_, err := ssmClient.PutParameter(in.Context, &ssm.PutParameterInput{
			Name:  paramName,
			Type:  param.Type,
			Value: paramValue,
		})
		if err != nil {
			return &errors.Object{
				Id:     "5b3c9f35-90ce-447e-a452-a4c07b3c3626",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to create parameter.",
				Cause:  err.Error(),
			}
		}
	case UpdateOptionUpdate:
		_, err := ssmClient.PutParameter(in.Context, &ssm.PutParameterInput{
			Name:      paramName,
			Type:      param.Type,
			Value:     paramValue,
			Overwrite: val.PtrRef(true),
		})
		if err != nil {
			return &errors.Object{
				Id:     "0b891402-f613-4fde-bdbe-923349ff6fc3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to update parameter.",
				Cause:  err.Error(),
			}
		}
	case UpdateOptionDelete:
		_, err := ssmClient.DeleteParameter(in.Context, &ssm.DeleteParameterInput{
			Name: paramName,
		})
		if err != nil {
			return &errors.Object{
				Id:     "7a9b05ac-dc36-4089-b392-f7c045ceaf47",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to delete parameter.",
				Cause:  err.Error(),
			}
		}
	}

	return nil
}
