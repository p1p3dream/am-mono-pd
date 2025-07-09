package conf

import (
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
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

	Casbin       *val.Cache[string, *casbin.Enforcer]
	Duration     *val.Cache[string, time.Duration]
	OpenSearch   *val.Cache[string, *opensearch.Client]
	Paseto       *val.Cache[string, *gconf.PasetoCacheItem]
	PGxPool      *val.Cache[string, *pgxpool.Pool]
	Valkey       *val.Cache[string, valkey.Client]
	ValkeyScript *val.Cache[string, *valkey.Lua]
}

type File struct {
	DeploymentEnvironment    int    `json:"-" yaml:"-"`
	DeploymentEnvironmentStr string `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`

	HttpServer *gconf.HttpServer `json:"http_server,omitempty" yaml:"http_server,omitempty"`

	Casbin     map[string]*gconf.Casbin     `json:"casbin,omitempty" yaml:"casbin,omitempty"`
	Duration   map[string]string            `json:"duration,omitempty" yaml:"duration,omitempty"`
	Flags      []string                     `json:"flags,omitempty" yaml:"flags,omitempty"`
	OpenSearch map[string]*gconf.OpenSearch `json:"opensearch,omitempty" yaml:"opensearch,omitempty"`
	Paseto     map[string]*gconf.Paseto     `json:"paseto,omitempty" yaml:"paseto,omitempty"`
	Postgres   map[string]*gconf.Postgres   `json:"postgres,omitempty" yaml:"postgres,omitempty"`
	Valkey     map[string]*gconf.Valkey     `json:"valkey,omitempty" yaml:"valkey,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`
}

func Resolve(configPath string) (*Config, error) {
	envKey := "ABODEMINE_SAAS_CONFIG_PATH"

	f := new(File)

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "48ffd9f3-3329-4dac-adad-c057aa70b78a")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "7d09a1ab-46b6-45d6-8e3d-e24874cb7c2b")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "ebba165b-571d-4f72-a4d5-b59d7cc012c5",
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
			return errors.Forward(err, "dcfdbb46-7f61-4bb4-a892-be653d7d0ba7")
		}
	}

	config.Casbin = val.NewCache[string, *casbin.Enforcer]()

	for k, v := range file.Casbin {
		item, err := gconf.LoadCasbin(v)
		if err != nil {
			return errors.Forward(err, "45496225-149c-4b02-88a8-6c6c345fe1b2")
		}

		config.Casbin.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Casbin configuration.")
	}

	config.Duration = val.NewCache[string, time.Duration]()

	for k, v := range file.Duration {
		item, err := gconf.LoadDuration(v)
		if err != nil {
			return errors.Forward(err, "1c7743de-ee53-4354-a86d-81fd111f286e")
		}

		config.Duration.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Duration configuration.")
	}

	if err := flags.ValidateMany(file.Flags); err != nil {
		return errors.Forward(err, "0a079b5d-6a07-49bd-9482-28c0a5675fc6")
	}

	for _, flag := range file.Flags {
		log.Info().Str("flag", flag).Msg("Loaded flag.")
	}

	config.OpenSearch = val.NewCache[string, *opensearch.Client]()

	for k, v := range file.OpenSearch {
		item, err := gconf.LoadOpenSearch(v)
		if err != nil {
			return errors.Forward(err, "e2431061-15f7-4c92-861b-e9be9bc79e52")
		}

		config.OpenSearch.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded OpenSearch configuration.")
	}

	config.Paseto = val.NewCache[string, *gconf.PasetoCacheItem]()

	for k, v := range file.Paseto {
		item, err := gconf.LoadPaseto(v)
		if err != nil {
			return errors.Forward(err, "f506a756-b19c-469d-8eff-6574a9a4b948")
		}

		config.Paseto.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Paseto configuration.")
	}

	config.PGxPool = val.NewCache[string, *pgxpool.Pool]()

	for k, v := range file.Postgres {
		item, err := gconf.LoadPostgres(v)
		if err != nil {
			return errors.Forward(err, "2f66a6f0-7ee8-48ff-8903-0e1530cefe58")
		}

		config.PGxPool.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded PGxPool configuration.")
	}

	config.Valkey = val.NewCache[string, valkey.Client]()
	config.ValkeyScript = val.NewCache[string, *valkey.Lua]()

	for k, v := range file.Valkey {
		item, err := gconf.LoadValkey(v)
		if err != nil {
			return errors.Forward(err, "88f6d8cf-7cd2-4924-8ad2-0ffa0b0fe6d4")
		}

		config.Valkey.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Valkey configuration.")

		for l, w := range v.Scripts {
			if config.ValkeyScript.Has(l) {
				return &errors.Object{
					Id:     "04f801b4-e364-4f33-93ee-fe26795b7f2a",
					Code:   errors.Code_ALREADY_EXISTS,
					Detail: "A script with the provided name has been registered already.",
				}
			}

			item, err := gconf.LoadValkeyScript(w)
			if err != nil {
				return errors.Forward(err, "25d332bd-e5c8-4a2a-9a27-28babbc98236")
			}

			config.ValkeyScript.Set(l, item)
			log.Info().Str("key", l).Msg("Loaded ValkeyScript configuration.")
		}
	}

	return nil
}

func ResolveAndLoad(path string) (*Config, error) {
	config, err := Resolve(path)
	if err != nil {
		return nil, errors.Forward(err, "0ca2c43b-1a49-499a-987b-fb08ade0b28c")
	}

	if err := Load(config); err != nil {
		return nil, errors.Forward(err, "a0821a7e-2137-4066-bc0a-1b2f8dba8dc2")
	}

	return config, nil
}

var configOnce sync.Once
var configOnceValue *Config

// MustResolveAndLoadOnce must be used only for tests.
func MustResolveAndLoadOnce() *Config {
	configOnce.Do(func() {
		config, err := ResolveAndLoad("")
		if err != nil {
			panic(err)
		}

		// LoadPostgresFixtureRunner(config.PGxPool.Get("default"))

		configOnceValue = config
	})

	return configOnceValue
}
