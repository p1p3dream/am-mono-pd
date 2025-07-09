package worker

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/stretchr/testify/assert"
	"github.com/valkey-io/valkey-go"

	"abodemine/domains/arc"
	"abodemine/lib/gconf"
)

type MockDomain struct {
	deploymentEnvironment int
}

func (d *MockDomain) CreateRequest(in *arc.CreateRequestInput) (*arc.Request, error) {
	return &arc.Request{}, nil
}

func (d *MockDomain) CreateServerSession(r *arc.Request, in *arc.CreateServerSessionInput) (arc.ServerSession, error) {
	return nil, nil
}

func (d *MockDomain) SelectServerSession(r *arc.Request, in *arc.SelectServerSessionInput) (arc.ServerSession, error) {
	return nil, nil
}

func (d *MockDomain) DeleteServerSession(r *arc.Request, in *arc.DeleteServerSessionInput) error {
	return nil
}

func (d *MockDomain) DeploymentEnvironment() int {
	return d.deploymentEnvironment
}

func (d *MockDomain) SelectAWS(k string) (aws.Config, error) {
	return aws.Config{}, nil
}

func (d *MockDomain) SelectCasbin(k string) (*casbin.Enforcer, error) {
	return nil, nil
}

func (d *MockDomain) SelectDuration(k string) (time.Duration, error) {
	return 0, nil
}

func (d *MockDomain) SelectOpenSearch(k string) (*opensearch.Client, error) {
	return nil, nil
}

func (d *MockDomain) SelectPaseto(k string) (*gconf.PasetoCacheItem, error) {
	return nil, nil
}

func (d *MockDomain) SelectPgxPool(k string) (*pgxpool.Pool, error) {
	return nil, nil
}

func (d *MockDomain) SelectValkey(k string) (valkey.Client, error) {
	return nil, nil
}

func (d *MockDomain) SelectValkeyScript(k string) (*valkey.Lua, error) {
	return nil, nil
}

type MockOpenSearchClient struct{}

func (m *MockOpenSearchClient) Bulk(ctx context.Context, body io.Reader, o ...func(*opensearchapi.BulkRequest)) (*opensearchapi.Response, error) {
	return &opensearchapi.Response{}, nil
}

func TestDomain_CreateRequest(t *testing.T) {
	mockDomain := &MockDomain{
		deploymentEnvironment: 1,
	}

	testCases := []struct {
		name    string
		input   *arc.CreateRequestInput
		wantErr bool
	}{
		{
			name: "success with default values",
			input: &arc.CreateRequestInput{
				Context: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "success with custom ID",
			input: &arc.CreateRequestInput{
				Id:      uuid.New(),
				Context: context.Background(),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := mockDomain.CreateRequest(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, req)
		})
	}
}
