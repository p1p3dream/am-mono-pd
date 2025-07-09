package arc

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/valkey-io/valkey-go"

	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/val"
)

type Domain interface {
	CreateRequest(in *CreateRequestInput) (*Request, error)

	CreateServerSession(r *Request, in *CreateServerSessionInput) (ServerSession, error)
	SelectServerSession(r *Request, in *SelectServerSessionInput) (ServerSession, error)
	DeleteServerSession(r *Request, in *DeleteServerSessionInput) error

	DeploymentEnvironment() int
	Values() *gconf.Values

	SelectAWS(k string) (aws.Config, error)
	SelectCasbin(k string) (*casbin.Enforcer, error)
	SelectDuration(k string) (time.Duration, error)
	SelectOpenSearch(k string) (*opensearch.Client, error)
	SelectPaseto(k string) (*gconf.PasetoCacheItem, error)
	SelectPgxPool(k string) (*pgxpool.Pool, error)
	SelectValkey(k string) (valkey.Client, error)
	SelectValkeyScript(k string) (*valkey.Lua, error)
}

type domain struct {
	deploymentEnvironment int
	flags                 []string
	values                *gconf.Values

	aws          *val.Cache[string, aws.Config]
	casbin       *val.Cache[string, *casbin.Enforcer]
	duration     *val.Cache[string, time.Duration]
	openSearch   *val.Cache[string, *opensearch.Client]
	paseto       *val.Cache[string, *gconf.PasetoCacheItem]
	pgxPool      *val.Cache[string, *pgxpool.Pool]
	valkey       *val.Cache[string, valkey.Client]
	valkeyScript *val.Cache[string, *valkey.Lua]
}

type NewDomainInput struct {
	DeploymentEnvironment int
	Flags                 []string
	Values                *gconf.Values

	AWS          *val.Cache[string, aws.Config]
	Casbin       *val.Cache[string, *casbin.Enforcer]
	Duration     *val.Cache[string, time.Duration]
	OpenSearch   *val.Cache[string, *opensearch.Client]
	Paseto       *val.Cache[string, *gconf.PasetoCacheItem]
	PgxPool      *val.Cache[string, *pgxpool.Pool]
	Valkey       *val.Cache[string, valkey.Client]
	ValkeyScript *val.Cache[string, *valkey.Lua]
}

func NewDomain(in *NewDomainInput) Domain {
	return &domain{
		deploymentEnvironment: in.DeploymentEnvironment,
		flags:                 in.Flags,
		values: val.Ternary(
			in.Values == nil,
			gconf.NewEmptyValues(),
			in.Values,
		),

		aws: val.Ternary(
			in.AWS == nil,
			val.NewCache[string, aws.Config](),
			in.AWS,
		),
		casbin: val.Ternary(
			in.Casbin == nil,
			val.NewCache[string, *casbin.Enforcer](),
			in.Casbin,
		),
		duration: val.Ternary(
			in.Duration == nil,
			val.NewCache[string, time.Duration](),
			in.Duration,
		),
		openSearch: val.Ternary(
			in.OpenSearch == nil,
			val.NewCache[string, *opensearch.Client](),
			in.OpenSearch,
		),
		paseto: val.Ternary(
			in.Paseto == nil,
			val.NewCache[string, *gconf.PasetoCacheItem](),
			in.Paseto,
		),
		pgxPool: val.Ternary(
			in.PgxPool == nil,
			val.NewCache[string, *pgxpool.Pool](),
			in.PgxPool,
		),
		valkey: val.Ternary(
			in.Valkey == nil,
			val.NewCache[string, valkey.Client](),
			in.Valkey,
		),
		valkeyScript: val.Ternary(
			in.ValkeyScript == nil,
			val.NewCache[string, *valkey.Lua](),
			in.ValkeyScript,
		),
	}
}

func (dom *domain) DeploymentEnvironment() int {
	return dom.deploymentEnvironment
}

func (dom *domain) SelectAWS(k string) (aws.Config, error) {
	v, ok := dom.aws.Select(k)
	if !ok {
		return aws.Config{}, &errors.Object{
			Id:     "72cb20f4-680f-4450-b3d6-269aa8a58235",
			Code:   errors.Code_INTERNAL,
			Detail: "AWS item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectCasbin(k string) (*casbin.Enforcer, error) {
	v, ok := dom.casbin.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "97183cb2-83eb-42f6-9137-e3092fd7d07e",
			Code:   errors.Code_INTERNAL,
			Detail: "Casbin item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectDuration(k string) (time.Duration, error) {
	v, ok := dom.duration.Select(k)
	if !ok {
		return 0, &errors.Object{
			Id:     "f36a8e24-c71e-4ea4-acde-c6823ccbe717",
			Code:   errors.Code_INTERNAL,
			Detail: "Duration item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectOpenSearch(k string) (*opensearch.Client, error) {
	v, ok := dom.openSearch.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "fbe72913-b4dd-45d4-b355-88f35d441de6",
			Code:   errors.Code_INTERNAL,
			Detail: "OpenSearch item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectPaseto(k string) (*gconf.PasetoCacheItem, error) {
	v, ok := dom.paseto.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "68617f09-380c-4d55-850b-74ce83126a7d",
			Code:   errors.Code_INTERNAL,
			Detail: "Paseto item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectPgxPool(k string) (*pgxpool.Pool, error) {
	v, ok := dom.pgxPool.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "d06ca9d6-4cf2-491a-89e0-64fc97f1ec09",
			Code:   errors.Code_INTERNAL,
			Detail: "PgxPool item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectValkey(k string) (valkey.Client, error) {
	v, ok := dom.valkey.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "1bcb4e9d-86f3-4ef2-9122-df6172925da9",
			Code:   errors.Code_INTERNAL,
			Detail: "Valkey item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) SelectValkeyScript(k string) (*valkey.Lua, error) {
	v, ok := dom.valkeyScript.Select(k)
	if !ok {
		return nil, &errors.Object{
			Id:     "d95d14fa-74ff-4bfd-8a21-d3b0a6726fdc",
			Code:   errors.Code_INTERNAL,
			Detail: "ValkeyScript item not found.",
			Meta: map[string]any{
				"key": k,
			},
		}
	}

	return v, nil
}

func (dom *domain) Values() *gconf.Values {
	return dom.values
}
