package arc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/lib/unique"
	"abodemine/lib/val"
	"abodemine/projects/ci/conf"
)

func TestDomain_CreateServerSession(t *testing.T) {
	uCache := val.NewCache[string, uuid.UUID]()
	sCache := val.NewCache[string, string]()

	testCases := []*struct {
		name string
		in   *CreateServerSessionInput
		out  ServerSession
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "51e73ab3-7ec5-4f93-82db-314f38f7deb7",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-organization-id",
			in:   &CreateServerSessionInput{},
			err: &errors.Object{
				Id:   "c2a27e96-04df-40d8-818c-d8febdb3fb48",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-user-id",
			in: &CreateServerSessionInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "84d2ad91-2b61-4823-a683-114520dfa7c4",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "unsupported-session-type",
			in: &CreateServerSessionInput{
				OrganizationId: uuid.New(),
				UserId:         uuid.New(),
			},
			err: &errors.Object{
				Id:   "a18fdfa3-8c83-4f0c-9436-8455f42459ba",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-ttl",
			in: &CreateServerSessionInput{
				OrganizationId: uuid.New(),
				UserId:         uuid.New(),
				SessionType:    SessionTypeApiServer,
			},
			err: &errors.Object{
				Id:   "961054c6-9364-40c3-ba03-f9062fd4c27f",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &CreateServerSessionInput{
				OrganizationId: uCache.SetGet("ok-OrganizationId", uuid.New()),
				UserId:         uCache.SetGet("ok-UserId", uuid.New()),
				RoleName:       sCache.SetGet("ok-RoleName", unique.AlphaNum(8)),
				Timezone:       sCache.SetGet("ok-Timezone", unique.AlphaNum(8)),
				SessionType:    SessionTypeApiServer,
				Username:       sCache.SetGet("ok-Username", unique.AlphaNum(8)),
				TTL:            time.Minute,
			},
			out: &serverSessionContainer{
				payload: serverSessionPayload{
					OrganizationId: uCache.Get("ok-OrganizationId"),
					UserId:         uCache.Get("ok-UserId"),
					RoleName:       sCache.Get("ok-RoleName"),
					Timezone:       sCache.Get("ok-Timezone"),
					SessionType:    SessionTypeApiServer,
					Username:       sCache.Get("ok-Username"),
				},
			},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)
	dom := NewDomain(&NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	wantUuidVersion := uuid.New().Version().String()

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := dom.CreateRequest(&CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			have, err := dom.CreateServerSession(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to CreateServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			// Ensure returned value is OK.

			want := tc.out

			assert.Equal(st, wantUuidVersion, have.Id().Version().String(), "ServerSession.Id[0] is not UUIDv4")
			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "ServerSession.OrganizationId[0] mismatch")
			assert.Equal(st, want.UserId(), have.UserId(), "ServerSession.UserId[0] mismatch")
			assert.Equal(st, want.RoleName(), have.RoleName(), "ServerSession.RoleName[0] mismatch")
			assert.Equal(st, want.Timezone(), have.Timezone(), "ServerSession.Timezone[0] mismatch")
			assert.Equal(st, want.SessionType(), have.SessionType(), "ServerSession.SessionType[0] mismatch")
			assert.Equal(st, want.Username(), have.Username(), "ServerSession.Username[0] mismatch")

			// Ensure key is present.

			have, err = dom.SelectServerSession(r, &SelectServerSessionInput{
				OrganizationId: have.OrganizationId(),
				Id:             have.Id(),
				SessionType:    have.SessionType(),
			})
			if err != nil {
				st.Fatalf("failed to SelectServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			assert.Equal(st, wantUuidVersion, have.Id().Version().String(), "ServerSession.Id[1] is not UUIDv4")
			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "ServerSession.OrganizationId[1] mismatch")
			assert.Equal(st, want.UserId(), have.UserId(), "ServerSession.UserId[1] mismatch")
			assert.Equal(st, want.RoleName(), have.RoleName(), "ServerSession.RoleName[1] mismatch")
			assert.Equal(st, want.Timezone(), have.Timezone(), "ServerSession.Timezone[1] mismatch")
			assert.Equal(st, want.SessionType(), have.SessionType(), "ServerSession.SessionType[1] mismatch")
			assert.Equal(st, want.Username(), have.Username(), "ServerSession.Username[1] mismatch")
		})
	}
}

func TestDomain_SelectServerSession(t *testing.T) {
	uCache := val.NewCache[string, uuid.UUID]()

	testCases := []*struct {
		name string
		in   *SelectServerSessionInput
		out  ServerSession
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "1f1a2c75-8088-4f7b-9a60-c9ad1afc9306",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-organization-id",
			in:   &SelectServerSessionInput{},
			err: &errors.Object{
				Id:   "5603dbee-573a-4362-b58b-88e54ea61d99",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-id",
			in: &SelectServerSessionInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "7f5ee5ce-afeb-4a24-827e-2c6ecba61e1c",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &SelectServerSessionInput{
				OrganizationId: uCache.SetGet("ok-OrganizationId", uuid.New()),
				SessionType:    SessionTypeApiServer,
			},
			out: &serverSessionContainer{
				payload: serverSessionPayload{
					OrganizationId: uCache.Get("ok-OrganizationId"),
					UserId:         uuid.New(),
					Timezone:       unique.AlphaNum(8),
					SessionType:    SessionTypeApiServer,
					Username:       unique.AlphaNum(8),
				},
			},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)
	dom := NewDomain(&NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := dom.CreateRequest(&CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			want := tc.out

			if tc.err == nil {
				created, err := dom.CreateServerSession(r, &CreateServerSessionInput{
					OrganizationId: tc.in.OrganizationId,
					UserId:         want.UserId(),
					SessionType:    want.SessionType(),
					Timezone:       want.Timezone(),
					Username:       want.Username(),
					TTL:            time.Minute,
				})
				if err != nil {
					st.Fatalf("failed to CreateServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				tc.in.Id = created.Id()
			}

			have, err := dom.SelectServerSession(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to SelectServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "ServerSession.OrganizationId mismatch")
			assert.Equal(st, want.UserId(), have.UserId(), "ServerSession.UserId mismatch")
			assert.Equal(st, want.Timezone(), have.Timezone(), "ServerSession.Timezone mismatch")
			assert.Equal(st, want.SessionType(), have.SessionType(), "ServerSession.SessionType mismatch")
			assert.Equal(st, want.Username(), have.Username(), "ServerSession.Username mismatch")
		})
	}
}

func TestDomain_DeleteServerSession(t *testing.T) {
	uCache := val.NewCache[string, uuid.UUID]()

	testCases := []*struct {
		name string
		in   *DeleteServerSessionInput
		out  ServerSession
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "c2976214-464e-4f96-98c8-eba4394fdeb4",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-organization-id",
			in:   &DeleteServerSessionInput{},
			err: &errors.Object{
				Id:   "95145a17-80fe-4deb-a4ae-26425d34d7c0",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-id",
			in: &DeleteServerSessionInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "0dad30a3-4c5e-404a-a000-098f8b293a3e",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &DeleteServerSessionInput{
				OrganizationId: uCache.SetGet("ok-OrganizationId", uuid.New()),
				SessionType:    SessionTypeSaasServer,
			},
			out: &serverSessionContainer{
				payload: serverSessionPayload{
					OrganizationId: uCache.Get("ok-OrganizationId"),
					UserId:         uuid.New(),
					Timezone:       unique.AlphaNum(8),
					SessionType:    SessionTypeSaasServer,
					Username:       unique.AlphaNum(8),
				},
			},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)
	dom := NewDomain(&NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := dom.CreateRequest(&CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			want := tc.out

			if tc.err == nil {
				created, err := dom.CreateServerSession(r, &CreateServerSessionInput{
					OrganizationId: tc.in.OrganizationId,
					UserId:         want.UserId(),
					SessionType:    want.SessionType(),
					Timezone:       want.Timezone(),
					Username:       want.Username(),
					TTL:            time.Minute,
				})
				if err != nil {
					st.Fatalf("failed to CreateServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				tc.in.Id = created.Id()
			}

			err = dom.DeleteServerSession(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to DeleteServerSession: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			// Ensure key is not present.

			_, err = dom.SelectServerSession(r, &SelectServerSessionInput{
				OrganizationId: tc.in.OrganizationId,
				Id:             tc.in.Id,
				SessionType:    want.SessionType(),
			})
			if err == nil {
				st.Fatalf("expected error, got nil")
			}

			assert.Equal(st, "e3e4ded9-9b1a-483c-a23b-60e743f95ab5", err.(*errors.Object).Id, "Error.Id mismatch")
		})
	}
}
