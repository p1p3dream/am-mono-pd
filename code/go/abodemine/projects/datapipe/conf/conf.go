package conf

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"

	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/gconf"
	"abodemine/lib/logging"
	"abodemine/lib/val"
)

type Config struct {
	File *File `json:"file,omitempty" yaml:"file,omitempty"`

	AWS          *val.Cache[string, aws.Config]
	OpenSearch   *val.Cache[string, *opensearch.Client]
	PgxPool      *val.Cache[string, *pgxpool.Pool]
	Valkey       *val.Cache[string, valkey.Client]
	ValkeyScript *val.Cache[string, *valkey.Lua]
}

type DistributedLocker struct {
	TableName string            `json:"table_name,omitempty" yaml:"table_name,omitempty"`
	Keys      map[string]string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

type File struct {
	DeploymentEnvironment    int    `json:"-" yaml:"-"`
	DeploymentEnvironmentStr string `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`

	Flags      []string                     `json:"flags,omitempty" yaml:"flags,omitempty"`
	OpenSearch map[string]*gconf.OpenSearch `json:"opensearch,omitempty" yaml:"opensearch,omitempty"`
	Postgres   map[string]*gconf.Postgres   `json:"postgres,omitempty" yaml:"postgres,omitempty"`
	Valkey     map[string]*gconf.Valkey     `json:"valkey,omitempty" yaml:"valkey,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`

	DistributedLockers map[string]*DistributedLocker `json:"distributed_lockers,omitempty" yaml:"distributed_lockers,omitempty"`
	FileBufferSize     int                           `json:"file_buffer_size,omitempty" yaml:"file_buffer_size,omitempty"`

	Lambdas *Lambdas `json:"lambdas,omitempty" yaml:"lambdas,omitempty"`
}

func Resolve(f *File, configPath string) (*Config, error) {
	envKey := "ABODEMINE_DATAPIPE_CONFIG_PATH"

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "ef09862b-3717-4d3a-a0bb-bfd33c6bd7d6")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "4fb4044a-cf37-4d0c-83ec-8bf1e1ed2c96")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "06438a85-f0e7-4d9c-86ab-363878045489",
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
			return errors.Forward(err, "9de38903-1a81-40d1-a09a-c438a6980cd5")
		}
	}

	if err := flags.ValidateMany(file.Flags); err != nil {
		return errors.Forward(err, "4a2acc2f-019c-4b96-a886-7916f4058d84")
	}

	for _, flag := range file.Flags {
		log.Info().Str("flag", flag).Msg("Loaded flag.")
	}

	config.OpenSearch = val.NewCache[string, *opensearch.Client]()

	for k, v := range file.OpenSearch {
		item, err := gconf.LoadOpenSearch(v)
		if err != nil {
			return errors.Forward(err, "39b12093-728e-423b-8ceb-c14dfa5b2f80")
		}

		config.OpenSearch.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded OpenSearch configuration")
	}

	config.PgxPool = val.NewCache[string, *pgxpool.Pool]()

	for k, v := range file.Postgres {
		item, err := gconf.LoadPostgres(v)
		if err != nil {
			return errors.Forward(err, "5aff03c7-3abe-46dc-9d70-e0c232407d28")
		}

		config.PgxPool.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded PgxPool configuration")
	}

	config.Valkey = val.NewCache[string, valkey.Client]()
	config.ValkeyScript = val.NewCache[string, *valkey.Lua]()

	for k, v := range file.Valkey {
		item, err := gconf.LoadValkey(v)
		if err != nil {
			return errors.Forward(err, "97dd4e4d-4ba1-4d18-9c90-37d5eb1fac8a")
		}

		config.Valkey.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Valkey configuration")

		for l, w := range v.Scripts {
			if config.ValkeyScript.Has(l) {
				return &errors.Object{
					Id:     "6fa1d56f-a6bd-48b6-ac8d-c26979fd2b80",
					Code:   errors.Code_ALREADY_EXISTS,
					Detail: "A script with the provided name has been registered already.",
				}
			}

			item, err := gconf.LoadValkeyScript(w)
			if err != nil {
				return errors.Forward(err, "6d4c1443-ae61-4643-a552-113abb0f33c3")
			}

			config.ValkeyScript.Set(l, item)
			log.Info().Str("key", l).Msg("Loaded ValkeyScript configuration")
		}
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
		return nil, errors.Forward(err, "83c87ff0-2463-42b8-aa93-9ba299afad20")
	}

	if err := Load(c); err != nil {
		return nil, errors.Forward(err, "85e693d1-7a3d-4ddc-bc3c-7b331e986167")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, &errors.Object{
			Id:     "c1f1cef0-b0af-4c5d-b3db-3071dd8385ea",
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
