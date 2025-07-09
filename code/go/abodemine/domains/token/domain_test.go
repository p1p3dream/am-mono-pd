package token

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/lib/val"
	"abodemine/projects/ci/conf"
)

func TestDomain_CreateToken(t *testing.T) {
	uuidCache := val.NewCache[string, uuid.UUID]()
	uint32Cache := val.NewCache[string, uint32]()

	testCases := []*struct {
		name string
		in   *CreateTokenInput
		out  Token
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "4bc74386-db1a-4925-b859-8791b777eea4",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "unsupported-token-type",
			in: &CreateTokenInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "d833ed02-ce0f-4d4c-a59b-e6d8d6c0768f",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-ttl",
			in: &CreateTokenInput{
				OrganizationId: uuid.New(),
				TokenType:      TypeOneTimePassword,
			},
			err: &errors.Object{
				Id:   "e8660186-fe03-40f4-a46f-5621bc20f1be",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &CreateTokenInput{
				TokenType: TypeTokenExchange,
				TTL:       time.Minute,
				Body:      val.ByteArray16ToSlice(uuidCache.SetGet("ok-Body", uuid.New())),
				Quota:     uint32Cache.SetGet("ok-Quota", must.GenRandomUint32()),
			},
			out: &tokenContainer{tokenPayload{
				TokenType: TypeTokenExchange,
				Body:      val.ByteArray16ToSlice(uuidCache.Get("ok-Body")),
				quota:     uint32Cache.Get("ok-Quota"),
			}},
		},
		{
			name: "ok-with-organization",
			in: &CreateTokenInput{
				OrganizationId: uuidCache.SetGet("ok-with-organization-OrganizationId", uuid.New()),
				TokenType:      TypeOneTimePassword,
				TTL:            time.Minute,
				Quota:          uint32Cache.SetGet("ok-with-organization-Quota", must.GenRandomUint32()),
			},
			out: &tokenContainer{tokenPayload{
				OrganizationId: uuidCache.Get("ok-with-organization-OrganizationId"),
				TokenType:      TypeOneTimePassword,
				quota:          uint32Cache.Get("ok-with-organization-Quota"),
			}},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)

	arcDomain := arc.NewDomain(&arc.NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	dom := NewDomain(&NewDomainInput{})
	wantUuidVersion := uuid.New().Version().String()

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDomain.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			createTokenOut, err := dom.CreateToken(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to CreateToken: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			have := createTokenOut.Token

			// Ensure returned value is OK.

			want := tc.out

			assert.Equal(st, wantUuidVersion, have.Id().Version().String(), "Token.Id[0] is not UUIDv4")
			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "Token.OrganizationId[0] mismatch")
			assert.Equal(st, want.TokenType(), have.TokenType(), "Token.TokenType[0] mismatch")
			assert.Equal(st, want.Body(), have.Body(), "Token.Body[0] mismatch")
			assert.Equal(st, want.Quota(), have.Quota(), "Token.Quota[0] mismatch")

			// // Ensure key is present.

			have, err = dom.SelectToken(r, &SelectTokenInput{
				OrganizationId: have.OrganizationId(),
				Id:             have.Id(),
				TokenType:      have.TokenType(),
				ReturnQuota:    true,
			})
			if err != nil {
				st.Fatalf("failed to SelectToken: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			assert.Equal(st, wantUuidVersion, have.Id().Version().String(), "Token.Id[1] is not UUIDv4")
			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "Token.OrganizationId[1] mismatch")
			assert.Equal(st, want.TokenType(), have.TokenType(), "Token.TokenType[1] mismatch")
			assert.Equal(st, want.Body(), have.Body(), "Token.Body[1] mismatch")
			assert.Equal(st, want.Quota(), have.Quota(), "Token.Quota[1] mismatch")
		})
	}
}

func TestDomain_SelectToken(t *testing.T) {
	uuidCache := val.NewCache[string, uuid.UUID]()
	uint32Cache := val.NewCache[string, uint32]()

	testCases := []*struct {
		name                string
		in                  *SelectTokenInput
		existingQuota       uint32
		createToken         bool
		checkTokenIsDeleted bool
		out                 Token
		err                 *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "83823559-7879-49b8-a6c7-2c553bf248bc",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-id",
			in: &SelectTokenInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "a8a80671-7163-4acb-a1b8-357740fffbf4",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "token-not-found",
			in: &SelectTokenInput{
				Id:        uuid.New(),
				TokenType: TypeOneTimePassword,
			},
			err: &errors.Object{
				Id:   "b52821f0-3127-4595-9927-70147f482472",
				Code: errors.Code_NOT_FOUND,
			},
		},
		{
			name: "insufficient-quota",
			in: &SelectTokenInput{
				TokenType:       TypeTokenExchange,
				QuotaDecreaseBy: 2,
				ReturnQuota:     true,
			},
			existingQuota: 1,
			createToken:   true,
			out: &tokenContainer{tokenPayload{
				TokenType: TypeTokenExchange,
			}},
			err: &errors.Object{
				Id:   "2c4dc6dd-6844-4110-b99e-49ba1cf14e5c",
				Code: errors.Code_RESOURCE_EXHAUSTED,
			},
		},
		{
			name: "ok",
			in: &SelectTokenInput{
				TokenType:       TypeOneTimePassword,
				QuotaDecreaseBy: uint32Cache.SetGet("ok-QuotaDecreaseBy", uint32(must.GenRandomInt32Range(0, 100))),
				ReturnQuota:     true,
			},
			existingQuota: uint32Cache.SetGet("ok-Quota", must.GenRandomUint32()+100),
			out: &tokenContainer{tokenPayload{
				TokenType: TypeOneTimePassword,
				Body:      val.ByteArray16ToSlice(uuid.New()),
				quota:     uint32Cache.Get("ok-Quota") - uint32Cache.Get("ok-QuotaDecreaseBy"),
			}},
		},
		{
			name: "ok-with-organization",
			in: &SelectTokenInput{
				OrganizationId: uuidCache.SetGet("ok-with-organization-OrganizationId", uuid.New()),
				TokenType:      TypeOneTimePassword,
			},
			existingQuota: uint32Cache.SetGet("ok-with-organization-Quota", must.GenRandomUint32()),
			out: &tokenContainer{tokenPayload{
				OrganizationId: uuidCache.Get("ok-with-organization-OrganizationId"),
				TokenType:      TypeOneTimePassword,
				Body:           val.ByteArray16ToSlice(uuid.New()),
				quota:          uint32Cache.Get("ok-with-organization-Quota"),
			}},
		},
		{
			name: "ok-delete-on-quota-exhausted",
			in: &SelectTokenInput{
				TokenType:       TypeOneTimePassword,
				QuotaDecreaseBy: 1,
				ReturnQuota:     true,
			},
			existingQuota: 1,
			out: &tokenContainer{tokenPayload{
				TokenType: TypeOneTimePassword,
				Body:      val.ByteArray16ToSlice(uuid.New()),
				quota:     0,
			}},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)

	arcDomain := arc.NewDomain(&arc.NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	dom := NewDomain(&NewDomainInput{})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDomain.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			want := tc.out

			if tc.err == nil || tc.createToken {
				createTokenOut, err := dom.CreateToken(r, &CreateTokenInput{
					OrganizationId: tc.in.OrganizationId,
					TokenType:      want.TokenType(),
					Body:           want.Body(),
					Quota:          tc.existingQuota,
					TTL:            time.Minute,
				})
				if err != nil {
					st.Fatalf("failed to CreateToken: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				created := createTokenOut.Token
				tc.in.Id = created.Id()
				want.(*tokenContainer).payload.Id = created.Id()
			}

			have, err := dom.SelectToken(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to SelectToken: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			assert.Equal(st, want.Id(), have.Id(), "Token.Id mismatch")
			assert.Equal(st, want.OrganizationId(), have.OrganizationId(), "Token.OrganizationId mismatch")
			assert.Equal(st, want.TokenType(), have.TokenType(), "Token.TokenType mismatch")
			assert.Equal(st, want.Body(), have.Body(), "Token.Body mismatch")

			if tc.in.ReturnQuota {
				assert.Equal(st, want.Quota(), have.Quota(), "Token.Quota mismatch")
			}

			if tc.checkTokenIsDeleted {
				_, err := dom.SelectToken(r, &SelectTokenInput{
					OrganizationId: have.OrganizationId(),
					Id:             have.Id(),
				})
				if err == nil {
					st.Fatal("Token is present.")
				}

				assert.Equal(st, "b52821f0-3127-4595-9927-70147f482472", err.(*errors.Object).Id, "Error.Id mismatch")
			}
		})
	}
}

func TestDomain_DeleteToken(t *testing.T) {
	uuidCache := val.NewCache[string, uuid.UUID]()

	testCases := []*struct {
		name string
		in   *DeleteTokenInput
		out  Token
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "641d5822-bd75-4289-b5f2-e8ca88f94363",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-id",
			in: &DeleteTokenInput{
				OrganizationId: uuid.New(),
			},
			err: &errors.Object{
				Id:   "4f2fb8d2-e499-414c-a87f-687d9a98607d",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &DeleteTokenInput{
				TokenType: TypeTokenExchange,
			},
			out: &tokenContainer{tokenPayload{
				TokenType: TypeTokenExchange,
				Body:      val.ByteArray16ToSlice(uuid.New()),
				quota:     must.GenRandomUint32(),
			}},
		},
		{
			name: "ok-with-organization",
			in: &DeleteTokenInput{
				OrganizationId: uuidCache.SetGet("ok-with-organization-OrganizationId", uuid.New()),
				TokenType:      TypeTokenExchange,
			},
			out: &tokenContainer{tokenPayload{
				OrganizationId: uuidCache.Get("ok-with-organization-OrganizationId"),
				TokenType:      TypeTokenExchange,
				Body:           val.ByteArray16ToSlice(uuid.New()),
				quota:          must.GenRandomUint32(),
			}},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)

	arcDomain := arc.NewDomain(&arc.NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	dom := NewDomain(&NewDomainInput{})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDomain.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "	"))
			}

			want := tc.out

			if tc.err == nil {
				createTokenOut, err := dom.CreateToken(r, &CreateTokenInput{
					OrganizationId: tc.in.OrganizationId,
					TokenType:      want.TokenType(),
					Body:           want.Body(),
					Quota:          want.Quota(),
					TTL:            time.Minute,
				})
				if err != nil {
					st.Fatalf("failed to CreateToken: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				created := createTokenOut.Token
				tc.in.Id = created.Id()
				want.(*tokenContainer).payload.Id = created.Id()
			}

			err = dom.DeleteToken(r, tc.in)
			if err != nil {
				if tc.err == nil {
					st.Fatalf("failed to DeleteToken: %s", must.MarshalJSONIndent(err, "", "	"))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")

				return
			}

			// Ensure key is not present.

			_, err = dom.SelectToken(r, &SelectTokenInput{
				OrganizationId: tc.in.OrganizationId,
				Id:             tc.in.Id,
				TokenType:      tc.in.TokenType,
			})
			if err == nil {
				st.Fatalf("Token is present.")
			}

			assert.Equal(st, "b52821f0-3127-4595-9927-70147f482472", err.(*errors.Object).Id, "Error.Id mismatch")
		})
	}
}
