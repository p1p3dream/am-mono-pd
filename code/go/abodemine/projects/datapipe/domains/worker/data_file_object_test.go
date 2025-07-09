package worker

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/entities"
)

func TestDomain_InsertDataFileObject(t *testing.T) {
	testCases := []*struct {
		name string
		in   *InsertDataFileObjectInput
		out  *InsertDataFileObjectOutput
		err  *errors.Object
	}{
		{
			name: "missing-input",
			err: &errors.Object{
				Id:   "9f01b149-c492-4344-8dec-a9766d6b6de8",
				Code: errors.Code_INTERNAL,
			},
		},
		{
			name: "missing-data-file-object",
			in:   &InsertDataFileObjectInput{},
			err: &errors.Object{
				Id:   "4cfe9efd-080d-49eb-be48-96162c8164ae",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-hash",
			in: &InsertDataFileObjectInput{
				Entity: &entities.DataFileObject{},
			},
			err: &errors.Object{
				Id:   "6f8acc2c-567b-4b44-86bd-20ba59ff2849",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-filetype",
			in: &InsertDataFileObjectInput{
				Entity: &entities.DataFileObject{
					Hash: val.ByteArray16ToSlice(uuid.New()),
				},
			},
			err: &errors.Object{
				Id:   "a6b72469-9df4-4127-b176-68ee86da731e",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-status",
			in: &InsertDataFileObjectInput{
				Entity: &entities.DataFileObject{
					Hash:     val.ByteArray16ToSlice(uuid.New()),
					FileType: 1,
				},
			},
			err: &errors.Object{
				Id:   "62208df3-5381-4843-94ca-22919beafc53",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-filedir",
			in: &InsertDataFileObjectInput{
				Entity: &entities.DataFileObject{
					Hash:     val.ByteArray16ToSlice(uuid.New()),
					FileType: 1,
					Status:   1,
				},
			},
			err: &errors.Object{
				Id:   "f642f63d-a526-417a-9ece-140c86ab50c9",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "missing-filename",
			in: &InsertDataFileObjectInput{
				Entity: &entities.DataFileObject{
					Hash:     val.ByteArray16ToSlice(uuid.New()),
					FileType: 1,
					Status:   1,
					FileDir:  ".",
				},
			},
			err: &errors.Object{
				Id:   "db74b226-98da-4082-837b-d6efcf57bf40",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "ok",
			in: &InsertDataFileObjectInput{
				ReturnEntity: true,
				Entity: &entities.DataFileObject{
					Hash:     val.ByteArray16ToSlice(uuid.New()),
					FileType: 1,
					Status:   1,
					FileDir:  ".",
					FileName: "test.txt",
				},
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

	addressDom := address.NewDomain(&address.NewDomainInput{})

	dom := NewDomain(&NewDomainInput{
		Config:        config,
		AddressDomain: addressDom,
	})

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDom.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("Failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "    "))
			}

			have, err := dom.InsertDataFileObject(r, tc.in)
			if err == nil && tc.err != nil {
				st.Fatalf("Expected error: %s", must.MarshalJSONIndent(tc.err, "", "    "))
			} else if err != nil {
				if tc.err == nil {
					st.Fatalf("Failed to InsertDataFileObject: %s", must.MarshalJSONIndent(err, "", "    "))
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
