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
	DeploymentEnvironment int `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`

	HttpServer *gconf.HttpServer `json:"http_server,omitempty" yaml:"http_server,omitempty"`

	OpenSearch map[string]*gconf.OpenSearch `json:"opensearch,omitempty" yaml:"opensearch,omitempty"`
	Postgres   map[string]*gconf.Postgres   `json:"postgres,omitempty" yaml:"postgres,omitempty"`
	Valkey     map[string]*gconf.Valkey     `json:"valkey,omitempty" yaml:"valkey,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`
}

func Resolve(f *File, configPath string) (*Config, error) {
	envKey := "ABODEMINE_CI_CONFIG_PATH"

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "f069ff05-19e3-4f37-9480-c271767aab49")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "2770bd98-c316-46c5-ac2a-1c4dbdf1ef5c")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "50408fe8-51ad-4875-b99c-e9c9e2c265bc",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing File configuration.",
			Path:   "/file",
		}
	}

	config.OpenSearch = val.NewCache[string, *opensearch.Client]()

	for k, v := range file.OpenSearch {
		item, err := gconf.LoadOpenSearch(v)
		if err != nil {
			return errors.Forward(err, "3abbf4b9-0031-4070-8487-e0c7c4812b29")
		}

		config.OpenSearch.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded OpenSearch configuration")
	}

	config.PgxPool = val.NewCache[string, *pgxpool.Pool]()

	for k, v := range file.Postgres {
		item, err := gconf.LoadPostgres(v)
		if err != nil {
			return errors.Forward(err, "e6fb79b5-593b-4887-85be-700dd6bf17fb")
		}

		config.PgxPool.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded PgxPool configuration")
	}

	config.Valkey = val.NewCache[string, valkey.Client]()
	config.ValkeyScript = val.NewCache[string, *valkey.Lua]()

	for k, v := range file.Valkey {
		item, err := gconf.LoadValkey(v)
		if err != nil {
			return errors.Forward(err, "d9466b81-05e9-4293-a04b-cd156cbf83d2")
		}

		config.Valkey.Set(k, item)
		log.Debug().Str("key", k).Msg("Loaded Valkey configuration")

		for l, w := range v.Scripts {
			if config.ValkeyScript.Has(l) {
				return &errors.Object{
					Id:     "c5a3f5a1-d608-437f-8c26-5a3a3c83fe84",
					Code:   errors.Code_ALREADY_EXISTS,
					Detail: "A script with the provided name has been registered already.",
				}
			}

			item, err := gconf.LoadValkeyScript(w)
			if err != nil {
				return errors.Forward(err, "16b3f36e-94be-49a8-809c-c4a42de0c534")
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
		return nil, errors.Forward(err, "3e59401c-0fc3-42b1-b61c-200f647d0ec8")
	}

	if err := Load(c); err != nil {
		return nil, errors.Forward(err, "03ad1a29-91b0-4cde-ad0c-8968bf15ce49")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, &errors.Object{
			Id:     "e5af6231-c8c1-4dd5-a304-fe8b41d871e6",
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
