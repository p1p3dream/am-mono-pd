package conf

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/gconf"
	"abodemine/lib/logging"
	"abodemine/lib/val"
)

type Config struct {
	File *File `json:"file,omitempty" yaml:"file,omitempty"`

	AWS *val.Cache[string, aws.Config]
}

type File struct {
	DeploymentEnvironment    int    `json:"-" yaml:"-"`
	DeploymentEnvironmentStr string `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`

	Flags []string `json:"flags,omitempty" yaml:"flags,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`

	Lambdas *Lambdas `json:"lambdas,omitempty" yaml:"lambdas,omitempty"`
}

func Resolve(f *File, configPath string) (*Config, error) {
	envKey := "ABODEMINE_PACKER_CONFIG_PATH"

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "1f59cf3d-d139-4a95-b9c2-2b6d69732355")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "c26be9c7-849c-4097-8bd6-a8bd678a9a13")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "c35b3f63-fad3-4cd0-9ec1-889f33294929",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing File configuration.",
			Path:   "/file",
		}
	}

	file.DeploymentEnvironment = gconf.DeploymentEnvironmentFromString(file.DeploymentEnvironmentStr)
	log.Info().Str("deployment_environment", gconf.DeploymentEnvironmentToString(file.DeploymentEnvironment)).Send()

	// Load first to ensure config errors go to Sentry.
	if file.Sentry != nil {
		if _, err := logging.InitSentry(file.Sentry, file.DeploymentEnvironment, file.NoLogColor); err != nil {
			return errors.Forward(err, "b04245e5-46e9-43f7-bca6-3336a56a8f16")
		}
	}

	if err := flags.ValidateMany(file.Flags); err != nil {
		return errors.Forward(err, "28879570-17aa-409d-80bc-69207b0e0b27")
	}

	for _, flag := range file.Flags {
		log.Info().Str("flag", flag).Msg("Loaded flag.")
	}

	return nil
}

func ResolveAndLoad(ctx context.Context, path string) (*Config, error) {
	var f *File

	// Check if there is a File in the context.
	v := ctx.Value(gconf.ConfigFileCtxKey("file"))

	if v == nil {
		f = new(File)
	} else {
		f = v.(*File)
	}

	c, err := Resolve(f, path)
	if err != nil {
		return nil, errors.Forward(err, "3c1bca47-6a85-4416-85f7-c9af12501155")
	}

	if err := Load(c); err != nil {
		return nil, errors.Forward(err, "289b52bb-2a9d-488a-97c0-fbadf2b11de5")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, &errors.Object{
			Id:     "c89b4065-29e9-4b2c-9a78-82f97bfc6577",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to load default AWS config.",
			Cause:  err.Error(),
		}
	}

	c.AWS = val.NewCache[string, aws.Config]()
	c.AWS.Set("default", awsConfig)

	return c, nil
}

var configOnce sync.Once
var configOnceValue *Config

// MustResolveAndLoadOnce must be used only for tests.
func MustResolveAndLoadOnce(ctx context.Context) *Config {
	configOnce.Do(func() {
		config, err := ResolveAndLoad(ctx, "")
		if err != nil {
			panic(err)
		}

		// LoadPostgresFixtureRunner(config.PGxPool.Get("default"))

		configOnceValue = config
	})

	return configOnceValue
}
