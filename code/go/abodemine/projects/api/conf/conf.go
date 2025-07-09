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

	Values *gconf.Values

	Casbin       *val.Cache[string, *casbin.Enforcer]
	Duration     *val.Cache[string, time.Duration]
	OpenSearch   *val.Cache[string, *opensearch.Client]
	Paseto       *val.Cache[string, *gconf.PasetoCacheItem]
	PGxPool      *val.Cache[string, *pgxpool.Pool]
	Valkey       *val.Cache[string, valkey.Client]
	ValkeyScript *val.Cache[string, *valkey.Lua]
}

type File struct {
	DeploymentEnvironment    int                 `json:"-" yaml:"-"`
	DeploymentEnvironmentStr string              `json:"deployment_environment,omitempty" yaml:"deployment_environment,omitempty"`
	Flags                    []string            `json:"flags,omitempty" yaml:"flags,omitempty"`
	Values                   *gconf.ConfigValues `json:"values,omitempty" yaml:"values,omitempty"`

	HttpServer *gconf.HttpServer `json:"http_server,omitempty" yaml:"http_server,omitempty"`

	Casbin     map[string]*gconf.Casbin     `json:"casbin,omitempty" yaml:"casbin,omitempty"`
	Duration   map[string]string            `json:"duration,omitempty" yaml:"duration,omitempty"`
	OpenSearch map[string]*gconf.OpenSearch `json:"opensearch,omitempty" yaml:"opensearch,omitempty"`
	Paseto     map[string]*gconf.Paseto     `json:"paseto,omitempty" yaml:"paseto,omitempty"`
	Postgres   map[string]*gconf.Postgres   `json:"postgres,omitempty" yaml:"postgres,omitempty"`
	Valkey     map[string]*gconf.Valkey     `json:"valkey,omitempty" yaml:"valkey,omitempty"`

	Sentry *gconf.Sentry `json:"sentry,omitempty" yaml:"sentry,omitempty"`

	LogLevel   string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	NoLogColor bool   `json:"no_log_color,omitempty" yaml:"no_log_color,omitempty"`
}

func Resolve(configPath string) (*Config, error) {
	envKey := "ABODEMINE_API_CONFIG_PATH"

	f := new(File)

	if err := gconf.ResolveConfig(f, configPath, envKey); err != nil {
		return nil, errors.Forward(err, "5b1a8983-20a5-4e65-94b5-e647ea01dc4e")
	}

	if err := gconf.LoadZerolog(f.LogLevel, f.NoLogColor); err != nil {
		return nil, errors.Forward(err, "80d20100-fbaf-4f9c-a61c-1d6139fa86c8")
	}

	return &Config{
		File: f,
	}, nil
}

func Load(config *Config) error {
	file := config.File

	if file == nil {
		return &errors.Object{
			Id:     "5217e4b2-1f59-417b-b425-52732cb71600",
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
			return errors.Forward(err, "dc32e994-aa59-4d8e-b313-8ceda908ab3a")
		}
	}

	if err := flags.ValidateMany(file.Flags); err != nil {
		return errors.Forward(err, "10763bdf-5535-4de3-ad9e-6df0e3318209")
	}

	for _, flag := range file.Flags {
		log.Info().Str("flag", flag).Msg("Loaded flag.")
	}

	values, err := gconf.LoadConfigValues(file.Values)
	if err != nil {
		return errors.Forward(err, "016597a0-ffca-4666-8c1d-3af69744b4fd")
	}

	config.Values = values
	config.Casbin = val.NewCache[string, *casbin.Enforcer]()

	for k, v := range file.Casbin {
		item, err := gconf.LoadCasbin(v)
		if err != nil {
			return errors.Forward(err, "8c56128e-9e70-40fb-91ca-e3a6b3d9bb8b")
		}

		config.Casbin.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Casbin configuration.")
	}

	config.Duration = val.NewCache[string, time.Duration]()

	for k, v := range file.Duration {
		item, err := gconf.LoadDuration(v)
		if err != nil {
			return errors.Forward(err, "3ac8207a-a517-4581-a96e-70b4ce699d2c")
		}

		config.Duration.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Duration configuration.")
	}

	config.OpenSearch = val.NewCache[string, *opensearch.Client]()

	for k, v := range file.OpenSearch {
		item, err := gconf.LoadOpenSearch(v)
		if err != nil {
			return errors.Forward(err, "3f73a381-e11c-4d4b-af92-587fd7e790b1")
		}

		config.OpenSearch.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded OpenSearch configuration.")
	}

	config.Paseto = val.NewCache[string, *gconf.PasetoCacheItem]()

	for k, v := range file.Paseto {
		item, err := gconf.LoadPaseto(v)
		if err != nil {
			return errors.Forward(err, "791bc631-0043-4460-8d31-c51884b64726")
		}

		config.Paseto.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Paseto configuration.")
	}

	config.PGxPool = val.NewCache[string, *pgxpool.Pool]()

	for k, v := range file.Postgres {
		item, err := gconf.LoadPostgres(v)
		if err != nil {
			return errors.Forward(err, "17d135aa-0f0d-4c0f-9278-4e51a03d559e")
		}

		config.PGxPool.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded PGxPool configuration.")
	}

	config.Valkey = val.NewCache[string, valkey.Client]()
	config.ValkeyScript = val.NewCache[string, *valkey.Lua]()

	for k, v := range file.Valkey {
		item, err := gconf.LoadValkey(v)
		if err != nil {
			return errors.Forward(err, "416302ad-f8dc-416e-a2b8-5b83945c76e1")
		}

		config.Valkey.Set(k, item)
		log.Info().Str("key", k).Msg("Loaded Valkey configuration.")

		for l, w := range v.Scripts {
			if config.ValkeyScript.Has(l) {
				return &errors.Object{
					Id:     "ea0de2e8-89dd-44fc-adbd-06f7381f42ed",
					Code:   errors.Code_ALREADY_EXISTS,
					Detail: "A script with the provided name has been registered already.",
				}
			}

			item, err := gconf.LoadValkeyScript(w)
			if err != nil {
				return errors.Forward(err, "bbdafbf1-9a70-40a1-8238-e143941903ec")
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
		return nil, errors.Forward(err, "64ddd423-21df-4024-aee2-9e1f7b4cbdbd")
	}

	if err := Load(config); err != nil {
		return nil, errors.Forward(err, "8a657ed3-25ca-4c53-adfe-bfc2c9489d82")
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
