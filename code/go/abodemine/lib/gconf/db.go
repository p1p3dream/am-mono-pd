package gconf

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	"github.com/valkey-io/valkey-go"

	"abodemine/lib/errors"
)

type OpenSearch struct {
	Options map[string]string `json:"options,omitempty" yaml:"options,omitempty"`

	Addresses            []string `json:"addresses,omitempty" yaml:"addresses,omitempty"`
	Username             string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password             string   `json:"password,omitempty" yaml:"password,omitempty"`
	CompressRequestBody  bool     `json:"compress_request_body,omitempty" yaml:"compress_request_body,omitempty"`
	DiscoverNodesOnStart bool     `json:"discover_nodes_on_start,omitempty" yaml:"discover_nodes_on_start,omitempty"`
	CollectionID         string   `json:"collection_id,omitempty" yaml:"collection_id,omitempty"`
}

func LoadOpenSearch(config *OpenSearch) (*opensearch.Client, error) {
	if len(config.Addresses) == 0 {
		return nil, &errors.Object{
			Id:     "00a9013b-21ff-40bc-8f14-f0a8523c0c17",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing url.",
		}
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses:            config.Addresses,
		Username:             config.Username,
		Password:             config.Password,
		CompressRequestBody:  config.CompressRequestBody,
		DiscoverNodesOnStart: config.DiscoverNodesOnStart,
	})
	if err != nil {
		return nil, &errors.Object{
			Id:     "8eb7369a-a721-408b-a86c-fc69be17393b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create client.",
			Cause:  err.Error(),
		}
	}

	// For serverless OpenSearch, we need to set the collection ID in the request headers
	if config.CollectionID != "" {
		client.Transport = &collectionTransport{
			base:         client.Transport,
			collectionID: config.CollectionID,
		}
	}

	return client, nil
}

// collectionTransport wraps the base transport to add the collection ID header
type collectionTransport struct {
	base         opensearchtransport.Interface
	collectionID string
}

func (t *collectionTransport) Perform(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Collection-ID", t.collectionID)
	return t.base.Perform(req)
}

type Postgres struct {
	Options map[string]string `json:"options,omitempty" yaml:"options,omitempty"`

	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Dbname   string `json:"dbname,omitempty" yaml:"dbname,omitempty"`
	SslMode  string `json:"ssl_mode,omitempty" yaml:"ssl_mode,omitempty"`
}

func LoadPostgres(config *Postgres) (*pgxpool.Pool, error) {
	config.Host = strings.TrimSpace(config.Host)
	config.User = strings.TrimSpace(config.User)
	config.Dbname = strings.TrimSpace(config.Dbname)

	connectionString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s",
		config.Host,
		config.Port,
		config.Dbname,
		config.User,
	)

	config.SslMode = strings.TrimSpace(config.SslMode)

	if config.SslMode != "" {
		connectionString += fmt.Sprintf(" sslmode=%s", config.SslMode)
	}

	if config.Password != "" {
		connectionString += fmt.Sprintf(" password=%s", config.Password)
	}

	var pingAfterConnect bool
	var driver string

	if len(config.Options) > 0 {
		pingAfterConnect = strings.ToLower(strings.TrimSpace(config.Options["pingAfterConnect"])) == "true"
		driver = strings.ToLower(strings.TrimSpace(config.Options["driver"]))
	}

	ctx := context.Background()

	switch driver {
	case "pgx":
		pool, err := pgxpool.New(ctx, connectionString)
		if err != nil {
			return nil, &errors.Object{
				Id:     "a76a0175-79f2-4920-b1d1-b5dcde74ce35",
				Code:   errors.Code_UNAVAILABLE,
				Detail: "Failed to open pgx connection.",
				Cause:  err.Error(),
			}
		}

		if pingAfterConnect {
			if err := pool.Ping(ctx); err != nil {
				return nil, &errors.Object{
					Id:     "04e434a7-0b37-4b3a-955e-2a63fa32747a",
					Code:   errors.Code_UNAVAILABLE,
					Detail: "Failed to ping pgx connection.",
					Cause:  err.Error(),
				}
			}
		}

		return pool, nil
	default:
		return nil, &errors.Object{
			Id:     "b9eafcbd-b0e7-4265-969b-187cc9aa1d9f",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unknown driver.",
		}
	}
}

type Valkey struct {
	Nodes    []*ValkeyNode            `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Scripts  map[string]*ValkeyScript `json:"scripts,omitempty" yaml:"scripts,omitempty"`
	Username string                   `json:"username,omitempty" yaml:"username,omitempty"`
	Password string                   `json:"password,omitempty" yaml:"password,omitempty"`
	Db       int                      `json:"db,omitempty" yaml:"db,omitempty"`
}

type ValkeyNode struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
}

type ValkeyScript struct {
	Body string `json:"body,omitempty" yaml:"body,omitempty"`
	File string `json:"file,omitempty" yaml:"file,omitempty"`
}

func LoadValkey(config *Valkey) (valkey.Client, error) {
	node := config.Nodes[0]

	opts := valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", node.Host, node.Port)},
	}

	config.Username = strings.TrimSpace(config.Username)
	if config.Username != "" {
		opts.Username = config.Username
	}

	config.Password = strings.TrimSpace(config.Password)
	if config.Password != "" {
		opts.Password = config.Password
	}

	opts.SelectDB = int(config.Db)

	cli, err := valkey.NewClient(opts)
	if err != nil {
		return nil, &errors.Object{
			Id:     "50a317d8-314b-4f44-86f5-d51fd4de3e5f",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create client.",
			Cause:  err.Error(),
		}
	}

	return cli, nil
}

func LoadValkeyScript(config *ValkeyScript) (*valkey.Lua, error) {
	if config == nil {
		return nil, &errors.Object{
			Id:     "716ad845-497a-4708-a748-2e2815c5c55a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing configuration.",
			Path:   "/",
		}
	}

	var body string

	if config.File == "" {
		body = config.Body
	} else {
		b, err := os.ReadFile(config.File)
		if err != nil {
			return nil, &errors.Object{
				Id:     "b9bd71ed-763c-4c30-9eb3-4f49912103d2",
				Code:   errors.Code_FAILED_PRECONDITION,
				Detail: "Failed to read file.",
				Path:   "/file",
				Cause:  err.Error(),
			}
		}

		body = string(b)
	}

	return valkey.NewLuaScript(body), nil
}
