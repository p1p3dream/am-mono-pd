package first_american

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/must"
	"abodemine/lib/storage"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/entities"
)

func TestDataSource_CreateDataFileEntry(t *testing.T) {
	stringCache := val.NewCache[string, string]()
	int64Cache := val.NewCache[string, int64]()

	testCases := []*struct {
		name string
		in   *entities.CreateDataFileEntryInput
		out  *entities.DataFileEntry
		err  *errors.Object
	}{
		{
			name: "nil-input",
			err: &errors.Object{
				Id:   "a84a2463-5000-4ec7-9646-1fbec9a93d95",
				Code: errors.Code_INTERNAL,
			},
		},
		{
			name: "missing-storage-object",
			in:   &entities.CreateDataFileEntryInput{},
			err: &errors.Object{
				Id:   "c5ad6b0e-bfee-4617-8fae-672c0a93e0a2",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "non-directory-input",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("non-directory-input-dir", uuid.New().String()),
					Name:        stringCache.SetGet("non-directory-input-name", uuid.New().String()),
					Size:        int64Cache.SetGet("non-directory-input-size", must.GenRandomInt64()),
					IsDirectory: false,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("non-directory-input-dir"),
					Name:        stringCache.Get("non-directory-input-name"),
					Size:        int64Cache.Get("non-directory-input-size"),
					IsDirectory: false,
				}),
				Ignore: true,
			},
		},
		{
			name: "address",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("address-dir", uuid.New().String()),
					Name:        stringCache.SetGet("address-name", "20240605_Address"),
					Size:        int64Cache.SetGet("address-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("address-dir"),
					Name:        stringCache.Get("address-name"),
					Size:        int64Cache.Get("address-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAddress,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeAddressPriority,
					20240605,
				},
			},
		},
		{
			name: "address-master",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("address-master-dir", uuid.New().String()),
					Name:        stringCache.SetGet("address-master-name", "20250326_AddressMaster"),
					Size:        int64Cache.SetGet("address-master-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("address-master-dir"),
					Name:        stringCache.Get("address-master-name"),
					Size:        int64Cache.Get("address-master-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAddress,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeAddressPriority,
					20250326,
				},
			},
		},
		{
			name: "annual-before-20250326",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("annual-before-20250326-dir", uuid.New().String()),
					Name:        stringCache.SetGet("annual-before-20250326-name", "20240101_Annual"),
					Size:        int64Cache.SetGet("annual-before-20250326-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("annual-before-20250326-dir"),
					Name:        stringCache.Get("annual-before-20250326-name"),
					Size:        int64Cache.Get("annual-before-20250326-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAssessorAnnual,
				Ignore:   true,
			},
		},
		{
			name: "annual",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("annual-dir", uuid.New().String()),
					Name:        stringCache.SetGet("annual-name", "20250326_Annual"),
					Size:        int64Cache.SetGet("annual-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("annual-dir"),
					Name:        stringCache.Get("annual-name"),
					Size:        int64Cache.Get("annual-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAssessorAnnual,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					20250326,
					DataFileTypeAssessorAnnualPriority,
				},
			},
		},
		{
			name: "annual-after-20250326",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("annual-after-20250326-dir", uuid.New().String()),
					Name:        stringCache.SetGet("annual-after-20250326-name", "20250403_Annual"),
					Size:        int64Cache.SetGet("annual-after-20250326-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("annual-after-20250326-dir"),
					Name:        stringCache.Get("annual-after-20250326-name"),
					Size:        int64Cache.Get("annual-after-20250326-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAssessorAnnual,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					20250403,
					DataFileTypeAssessorAnnualPriority,
				},
			},
		},
		{
			name: "asr",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("asr-dir", uuid.New().String()),
					Name:        stringCache.SetGet("asr-name", "20240101_ASR"),
					Size:        int64Cache.SetGet("asr-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("asr-dir"),
					Name:        stringCache.Get("asr-name"),
					Size:        int64Cache.Get("asr-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeASR,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeASRPriority,
					20240101,
				},
			},
		},
		{
			name: "avm-power-before-2025",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("avm-power-before-2025-dir", uuid.New().String()),
					Name:        stringCache.SetGet("avm-power-before-2025-name", "20240326_AVMPower"),
					Size:        int64Cache.SetGet("avm-power-before-2025-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("avm-power-before-2025-dir"),
					Name:        stringCache.Get("avm-power-before-2025-name"),
					Size:        int64Cache.Get("avm-power-before-2025-size"),
					IsDirectory: true,
				}),
				Ignore: true,
			},
		},
		{
			name: "avm-power",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("avm-power-dir", uuid.New().String()),
					Name:        stringCache.SetGet("avm-power-name", "20250326_AVMPower"),
					Size:        int64Cache.SetGet("avm-power-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("avm-power-dir"),
					Name:        stringCache.Get("avm-power-name"),
					Size:        int64Cache.Get("avm-power-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAVMPower,
				Priorities: []int32{
					DataFileTypeAVMPowerPriority,
					20250326,
				},
			},
		},
		{
			name: "deed",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("deed-dir", uuid.New().String()),
					Name:        stringCache.SetGet("deed-name", "20240605_Deed"),
					Size:        int64Cache.SetGet("deed-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("deed-dir"),
					Name:        stringCache.Get("deed-name"),
					Size:        int64Cache.Get("deed-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeDeedMtg,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeDeedMtgPriority,
					20240605,
				},
			},
		},
		{
			name: "deed-mtg",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("deed-mtg-dir", uuid.New().String()),
					Name:        stringCache.SetGet("deed-mtg-name", "20250404_DeedMtg"),
					Size:        int64Cache.SetGet("deed-mtg-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("deed-mtg-dir"),
					Name:        stringCache.Get("deed-mtg-name"),
					Size:        int64Cache.Get("deed-mtg-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeDeedMtg,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeDeedMtgPriority,
					20250404,
				},
			},
		},
		{
			name: "hoa",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("hoa-dir", uuid.New().String()),
					Name:        stringCache.SetGet("hoa-name", "20240101_HOA"),
					Size:        int64Cache.SetGet("hoa-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("hoa-dir"),
					Name:        stringCache.Get("hoa-name"),
					Size:        int64Cache.Get("hoa-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeHOA,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeHOAPriority,
					20240101,
				},
			},
		},
		{
			name: "hoa-lien",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("hoa-lien-dir", uuid.New().String()),
					Name:        stringCache.SetGet("hoa-lien-name", "20240102_HOALien"),
					Size:        int64Cache.SetGet("hoa-lien-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("hoa-lien-dir"),
					Name:        stringCache.Get("hoa-lien-name"),
					Size:        int64Cache.Get("hoa-lien-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeHOALien,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeHOALienPriority,
					20240102,
				},
			},
		},
		{
			name: "hpi",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("hpi-dir", uuid.New().String()),
					Name:        stringCache.SetGet("hpi-name", "20250325_HPI"),
					Size:        int64Cache.SetGet("hpi-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("hpi-dir"),
					Name:        stringCache.Get("hpi-name"),
					Size:        int64Cache.Get("hpi-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeHPI,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeHPIPriority,
					20250325,
				},
			},
		},
		{
			name: "invl",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("invl-dir", uuid.New().String()),
					Name:        stringCache.SetGet("invl-name", "20240604_INVL"),
					Size:        int64Cache.SetGet("invl-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("invl-dir"),
					Name:        stringCache.Get("invl-name"),
					Size:        int64Cache.Get("invl-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeInvLien,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeInvLienPriority,
					20240604,
				},
			},
		},
		{
			name: "inv-lien",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("inv-lien-dir", uuid.New().String()),
					Name:        stringCache.SetGet("inv-lien-name", "20250407_InvLien"),
					Size:        int64Cache.SetGet("inv-lien-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("inv-lien-dir"),
					Name:        stringCache.Get("inv-lien-name"),
					Size:        int64Cache.Get("inv-lien-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeInvLien,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeInvLienPriority,
					20250407,
				},
			},
		},
		{
			name: "listing",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("listing-dir", uuid.New().String()),
					Name:        stringCache.SetGet("listing-name", "20240604_Listing"),
					Size:        int64Cache.SetGet("listing-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("listing-dir"),
					Name:        stringCache.Get("listing-name"),
					Size:        int64Cache.Get("listing-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeListing,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeListingPriority,
					20240604,
				},
			},
		},
		{
			name: "listings",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("listings-dir", uuid.New().String()),
					Name:        stringCache.SetGet("listings-name", "20250404_Listings"),
					Size:        int64Cache.SetGet("listings-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("listings-dir"),
					Name:        stringCache.Get("listings-name"),
					Size:        int64Cache.Get("listings-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeListing,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeListingPriority,
					20250404,
				},
			},
		},
		{
			name: "nod",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("nod-dir", uuid.New().String()),
					Name:        stringCache.SetGet("nod-name", "20240604_NOD"),
					Size:        int64Cache.SetGet("nod-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("nod-dir"),
					Name:        stringCache.Get("nod-name"),
					Size:        int64Cache.Get("nod-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeNOD,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeNODPriority,
					20240604,
				},
			},
		},
		{
			name: "pfc",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("pfc-dir", uuid.New().String()),
					Name:        stringCache.SetGet("pfc-name", "20250403_PFC"),
					Size:        int64Cache.SetGet("pfc-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("pfc-dir"),
					Name:        stringCache.Get("pfc-name"),
					Size:        int64Cache.Get("pfc-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypePFC,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypePFCPriority,
					20250403,
				},
			},
		},
		{
			name: "power-avmhist",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("power-avmhist-dir", uuid.New().String()),
					Name:        stringCache.SetGet("power-avmhist-name", "20250317_Power_AVMHIST_2019"),
					Size:        int64Cache.SetGet("power-avmhist-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("power-avmhist-dir"),
					Name:        stringCache.Get("power-avmhist-name"),
					Size:        int64Cache.Get("power-avmhist-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAVMPower,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeAVMPowerHistoryPriority,
					2019,
					20250317,
				},
			},
		},
		{
			name: "shape",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("shape-dir", uuid.New().String()),
					Name:        stringCache.SetGet("shape-name", "20250325_Shape"),
					Size:        int64Cache.SetGet("shape-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("shape-dir"),
					Name:        stringCache.Get("shape-name"),
					Size:        int64Cache.Get("shape-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeShape,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeShapePriority,
					20250325,
				},
			},
		},
		{
			name: "tax-hist",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("tax-hist-dir", uuid.New().String()),
					Name:        stringCache.SetGet("tax-hist-name", "20240604_TaxHist"),
					Size:        int64Cache.SetGet("tax-hist-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("tax-hist-dir"),
					Name:        stringCache.Get("tax-hist-name"),
					Size:        int64Cache.Get("tax-hist-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeTaxHistory,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeTaxHistoryPriority,
					20240604,
				},
			},
		},
		{
			name: "tax-history",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("tax-history-dir", uuid.New().String()),
					Name:        stringCache.SetGet("tax-history-name", "20250325_TaxHistory"),
					Size:        int64Cache.SetGet("tax-history-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("tax-history-dir"),
					Name:        stringCache.Get("tax-history-name"),
					Size:        int64Cache.Get("tax-history-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeTaxHistory,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeTaxHistoryPriority,
					20250325,
				},
			},
		},
		{
			name: "update-before-20250403",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("update-before-20250403-dir", uuid.New().String()),
					Name:        stringCache.SetGet("update-before-20250403-name", "20250325_Update"),
					Size:        int64Cache.SetGet("update-before-20250403-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("update-before-20250403-dir"),
					Name:        stringCache.Get("update-before-20250403-name"),
					Size:        int64Cache.Get("update-before-20250403-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAssessorUpdate,
				Ignore:   true,
			},
		},
		{
			name: "update",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("update-dir", uuid.New().String()),
					Name:        stringCache.SetGet("update-name", "20250403_Update"),
					Size:        int64Cache.SetGet("update-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("update-dir"),
					Name:        stringCache.Get("update-name"),
					Size:        int64Cache.Get("update-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeAssessorUpdate,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					20250403,
					DataFileTypeAssessorUpdatePriority,
				},
			},
		},
		{
			name: "val-hist",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("val-hist-dir", uuid.New().String()),
					Name:        stringCache.SetGet("val-hist-name", "20240604_ValHist"),
					Size:        int64Cache.SetGet("val-hist-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("val-hist-dir"),
					Name:        stringCache.Get("val-hist-name"),
					Size:        int64Cache.Get("val-hist-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeValueHistory,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeValueHistoryPriority,
					20240604,
				},
			},
		},
		{
			name: "value-hist",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("value-hist-dir", uuid.New().String()),
					Name:        stringCache.SetGet("value-hist-name", "20250325_ValueHist"),
					Size:        int64Cache.SetGet("value-hist-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("value-hist-dir"),
					Name:        stringCache.Get("value-hist-name"),
					Size:        int64Cache.Get("value-hist-size"),
					IsDirectory: true,
				}),
				FileType: DataFileTypeValueHistory,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeValueHistoryPriority,
					20250325,
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

	dataSource := &DataSource{}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.name), func(st *testing.T) {
			r, err := arcDom.CreateRequest(&arc.CreateRequestInput{
				Context: ctx,
			})
			if err != nil {
				st.Fatalf("Failed to CreateRequest: %s", must.MarshalJSONIndent(err, "", "    "))
			}

			have, err := dataSource.CreateDataFileEntry(r, tc.in)
			if err == nil && tc.err != nil {
				st.Fatalf("Expected error: %s", must.MarshalJSONIndent(tc.err, "", "    "))
			} else if err != nil {
				if tc.err == nil {
					st.Fatalf("Failed to CreateDataFileEntry: %s", must.MarshalJSONIndent(err, "", "    "))
				}

				want := tc.err
				have := err.(*errors.Object)

				assert.Equal(st, want.Id, have.Id, "Error.Id mismatch")
				assert.Equal(st, want.Code, have.Code, "Error.Code mismatch")
				assert.Equal(st, want.Label, have.Label, "Error.Label mismatch")

				return
			}

			want := tc.out

			assert.Equal(st, want.StorageObject.Dir, have.StorageObject.Dir, "StorageObject.Dir mismatch")
			assert.Equal(st, want.StorageObject.Name, have.StorageObject.Name, "StorageObject.Name mismatch")
			assert.Equal(st, want.StorageObject.Size, have.StorageObject.Size, "StorageObject.Size mismatch")
			assert.Equal(st, want.StorageObject.IsDirectory(), have.StorageObject.IsDirectory(), "StorageObject.IsDirectory() mismatch")
			assert.Equal(st, want.EnterDirectory, have.EnterDirectory, "EnterDirectory mismatch")
			assert.Equal(st, want.FileType, have.FileType, "FileType mismatch")
			assert.Equal(st, want.Ignore, have.Ignore, "Ignore mismatch")
			assert.Equal(st, want.IgnoreSubDirs, have.IgnoreSubDirs, "IgnoreSubDirs mismatch")
			assert.Equal(st, want.Priorities, have.Priorities, "Priorities mismatch")
		})
	}
}
