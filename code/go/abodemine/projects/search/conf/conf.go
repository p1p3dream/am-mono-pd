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

type File struct {
	DeploymentEnvironment    int               `json:"-" yaml:"-"`
	DeploymentEnvironmentStr string            `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`
	HttpServer               *gconf.HttpServer `json:"http_server,omitempty" yaml:"http_server,omitempty"`

	OpenSearch map[string]*gconf.OpenSearch `json:"opensearch,omitempty" yaml:"opensearch,omitempty"`
	Postgres   map[string]*gconf.Postgres   `json:"postgres,omitempty" yaml:"postgres,omitempty"`
	Valkey     map[string]*gconf.Valkey     `json:"valkey,omitempty" yaml:"valkey,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`
}

func Resolve(f *File, configPath string) (*Config, error) {
	envKey := "ABODEMINE_SEARCH_CONFIG_PATH"

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "11b7ada8-1cad-4ca4-88a1-ec9fee2ffb46")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "246c4fa6-9f31-4b5f-9bde-b6bb9e90c3f9")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "c458bf35-9caa-471a-ad64-bce22001966d",
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
			return errors.Forward(err, "dfc4faed-4002-469a-a9c0-856604b2380c")
		}
	}

	config.OpenSearch = val.NewCache[string, *opensearch.Client]()

	for k, v := range file.OpenSearch {
		item, err := gconf.LoadOpenSearch(v)
		if err != nil {
			return errors.Forward(err, "5883b2b0-1b35-4c0d-83c6-923da449b2a6")
		}

		config.OpenSearch.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded OpenSearch configuration")
	}

	config.PgxPool = val.NewCache[string, *pgxpool.Pool]()

	for k, v := range file.Postgres {
		item, err := gconf.LoadPostgres(v)
		if err != nil {
			return errors.Forward(err, "c618ae16-3fba-49d7-9212-0da0629afaa8")
		}

		config.PgxPool.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded PgxPool configuration")
	}

	config.Valkey = val.NewCache[string, valkey.Client]()
	config.ValkeyScript = val.NewCache[string, *valkey.Lua]()

	for k, v := range file.Valkey {
		item, err := gconf.LoadValkey(v)
		if err != nil {
			return errors.Forward(err, "d1d4ce75-726d-4985-a629-2d7b6dbe7134")
		}

		config.Valkey.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded Valkey configuration")

		for l, w := range v.Scripts {
			if config.ValkeyScript.Has(l) {
				return &errors.Object{
					Id:     "4725ba87-a3da-4bb1-93c9-7942c2b47945",
					Code:   errors.Code_ALREADY_EXISTS,
					Detail: "A script with the provided name has been registered already.",
				}
			}

			item, err := gconf.LoadValkeyScript(w)
			if err != nil {
				return errors.Forward(err, "edc9a0b3-6bcb-4aab-9697-94a673de0dc6")
			}

			config.ValkeyScript.Set(l, item)
			log.Debug().Str("key", l).Msg("Loaded ValkeyScript configuration")
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
		return nil, errors.Forward(err, "9ee6cabe-ca6e-4020-9c02-a2cb2751354f")
	}

	if err := Load(c); err != nil {
		return nil, errors.Forward(err, "c9c9e7ff-6819-4e77-8cb0-a15f9491f466")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, &errors.Object{
			Id:     "320b3aec-e44b-430f-b0b7-bb007c866efa",
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
