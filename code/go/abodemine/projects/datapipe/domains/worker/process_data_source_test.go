package worker

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/lib/ptr"
	"abodemine/lib/storage"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/domains/partners/attom_data"
	"abodemine/projects/datapipe/domains/partners/first_american"
)

func TestDomain_ProcessDataSourceDirV2(t *testing.T) {
	uuidCache := val.NewCache[string, uuid.UUID]()

	testCases := []*struct {
		name string
		in   *ProcessDataSourceDirInput
		out  *ProcessDataSourceDirOutput
		err  *errors.Object
	}{
		{
			name: "ok-attom-data",
			in: &ProcessDataSourceDirInput{
				Backend: &storage.LocalBackend{
					FilesystemPath: filepath.Join(
						envAbodemineWorkspace,
						"assets/projects/datapipe/test-data/attom-data",
					),
				},
				DataSource:     &attom_data.DataSource{},
				FileBufferSize: 3,
				IsRootDir:      true,
				Meta: map[string]any{
					"batchId": uuidCache.SetGet("ok-attom-data-batchId", uuid.New()),
				},
				PartnerId: attom_data.PartnerId,
				Path:      "ftp",
				WorkerId:  ptr.UUID(uuidCache.SetGet("ok-attom-data-workerId", uuid.New())),
			},
		},
		{
			name: "ok-first-american",
			in: &ProcessDataSourceDirInput{
				Backend: &storage.LocalBackend{
					FilesystemPath: filepath.Join(
						envAbodemineWorkspace,
						"assets/projects/datapipe/test-data/first-american",
					),
				},
				DataSource:     &first_american.DataSource{},
				FileBufferSize: 3,
				IsRootDir:      true,
				Meta: map[string]any{
					"batchId": uuidCache.SetGet("ok-first-american-batchId", uuid.New()),
				},
				PartnerId: first_american.PartnerId,
				Path:      "ftp",
				WorkerId:  ptr.UUID(uuidCache.SetGet("ok-first-american-workerId", uuid.New())),
			},
		},
	}

	ctx := context.Background()
	config := conf.MustResolveAndLoadOnce(ctx)
	arcDom := arc.NewDomain(&arc.NewDomainInput{
		AWS:          config.AWS,
		OpenSearch:   config.OpenSearch,
		PgxPool:      config.PgxPool,
		Valkey:       config.Valkey,
		ValkeyScript: config.ValkeyScript,
	})

	dom := NewDomain(&NewDomainInput{
		Config: config,
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDom.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("Failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "    "))
			}

			have, err := dom.ProcessDataSourceDirV2(r, tc.in)
			if err == nil && tc.err != nil {
				st.Fatalf("Expected error: %s", must.MarshalJSONIndent(tc.err, "", "    "))
			} else if err != nil {
				if tc.err == nil {
					st.Fatalf("Failed to ProcessDataSourceDirV2: %s", must.MarshalJSONIndent(err, "", "    "))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")
				assert.Equal(st, want.Label, have.Label, "Error.Label mismatch")

				return
			}

			_ = have
		})
	}
}
