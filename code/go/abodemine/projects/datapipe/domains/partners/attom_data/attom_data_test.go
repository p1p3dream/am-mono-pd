package attom_data

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
				Id:   "5e376d72-7d2d-4a0b-a731-421d45ba5b75",
				Code: errors.Code_INTERNAL,
			},
		},
		{
			name: "missing-storage-object",
			in:   &entities.CreateDataFileEntryInput{},
			err: &errors.Object{
				Id:   "2f63170e-74a3-4cb6-b5ba-7f7519138634",
				Code: errors.Code_INVALID_ARGUMENT,
			},
		},
		{
			name: "refresh-dir",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.SetGet("refresh-dir-input-dir", "ftp"),
					Name:        stringCache.SetGet("refresh-dir-input-name", "Refresh"),
					Size:        int64Cache.SetGet("refresh-dir-input-size", must.GenRandomInt64()),
					IsDirectory: true,
				}),
				Path: stringCache.SetGet("refresh-dir-input-path", "/ftp/Refresh"),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:         stringCache.Get("refresh-dir-input-dir"),
					Name:        stringCache.Get("refresh-dir-input-name"),
					Size:        int64Cache.Get("refresh-dir-input-size"),
					IsDirectory: true,
				}),
				FileType:       entities.DataFileTypeSelectedDirectory,
				EnterDirectory: true,
				IgnoreSubDirs:  true,
			},
		},
		{
			name: "assignment-release",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("assignment-release-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("assignment-release-input-name", "RENTNECTAR_ASSIGNMENTRELEASE_0091.zip"),
					Size: int64Cache.SetGet("assignment-release-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("assignment-release-input-dir"),
					Name: stringCache.Get("assignment-release-input-name"),
					Size: int64Cache.Get("assignment-release-input-size"),
				}),
				FileType: DataFileTypeAssignmentRelease,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeAssignmentReleasePriority,
					91,
				},
			},
		},
		{
			name: "avm",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("avm-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("avm-input-name", "RENTNECTAR_AVM_0019.zip"),
					Size: int64Cache.SetGet("avm-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("avm-input-dir"),
					Name: stringCache.Get("avm-input-name"),
					Size: int64Cache.Get("avm-input-size"),
				}),
				FileType: DataFileTypeAvm,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeAvmPriority,
					19,
				},
			},
		},
		{
			name: "building-permit",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("building-permit-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("building-permit-input-name", "RENTNECTAR_BUILDINGPERMIT_0024.zip"),
					Size: int64Cache.SetGet("building-permit-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("building-permit-input-dir"),
					Name: stringCache.Get("building-permit-input-name"),
					Size: int64Cache.Get("building-permit-input-size"),
				}),
				FileType: DataFileTypeBuildingPermit,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeBuildingPermitPriority,
					24,
				},
			},
		},
		{
			name: "building-permit-classifiers",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("building-permit-classifiers-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("building-permit-classifiers-input-name", "RENTNECTAR_BUILDINGPERMITCLASSIFIERS_0021.zip"),
					Size: int64Cache.SetGet("building-permit-classifiers-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("building-permit-classifiers-input-dir"),
					Name: stringCache.Get("building-permit-classifiers-input-name"),
					Size: int64Cache.Get("building-permit-classifiers-input-size"),
				}),
				FileType: DataFileTypeBuildingPermitClassifiers,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeBuildingPermitClassifiersPriority,
					21,
				},
			},
		},
		{
			name: "building-permit-delete",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("building-permit-delete-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("building-permit-delete-input-name", "RENTNECTAR_BUILDINGPERMITDELETE_0022.zip"),
					Size: int64Cache.SetGet("building-permit-delete-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("building-permit-delete-input-dir"),
					Name: stringCache.Get("building-permit-delete-input-name"),
					Size: int64Cache.Get("building-permit-delete-input-size"),
				}),
				FileType: DataFileTypeBuildingPermitDelete,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeBuildingPermitDeletePriority,
					22,
				},
			},
		},
		{
			name: "building-permit-status",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("building-permit-status-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("building-permit-status-input-name", "RENTNECTAR_BUILDINGPERMITSTATUS_0023.zip"),
					Size: int64Cache.SetGet("building-permit-status-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("building-permit-status-input-dir"),
					Name: stringCache.Get("building-permit-status-input-name"),
					Size: int64Cache.Get("building-permit-status-input-size"),
				}),
				FileType: DataFileTypeBuildingPermitStatus,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeBuildingPermitStatusPriority,
					23,
				},
			},
		},
		{
			name: "cf",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("cf-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("cf-input-name", "Rentnectar_CF_20240902.zip"),
					Size: int64Cache.SetGet("cf-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("cf-input-dir"),
					Name: stringCache.Get("cf-input-name"),
					Size: int64Cache.Get("cf-input-size"),
				}),
				FileType: DataFileTypeCF,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeCFPriority,
					20240902,
				},
			},
		},
		{
			name: "daily-foreclosure",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("daily-foreclosure-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("daily-foreclosure-input-name", "RENTNECTAR_DAILY_FORECLOSURE_0091.zip"),
					Size: int64Cache.SetGet("daily-foreclosure-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("daily-foreclosure-input-dir"),
					Name: stringCache.Get("daily-foreclosure-input-name"),
					Size: int64Cache.Get("daily-foreclosure-input-size"),
				}),
				FileType: DataFileTypeDailyForeclosure,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeDailyForeclosurePriorityGroup,
					DataFileTypeDailyForeclosureRegularPriorityGroup,
					91,
					DataFileTypeDailyForeclosurePriority,
				},
			},
		},
		{
			name: "daily-foreclosure-refresh",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("daily-foreclosure-refresh-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("daily-foreclosure-refresh-input-name", "RENTNECTAR_REFRESH_FORECLOSURE_0002_001.zip"),
					Size: int64Cache.SetGet("daily-foreclosure-refresh-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("daily-foreclosure-refresh-input-dir"),
					Name: stringCache.Get("daily-foreclosure-refresh-input-name"),
					Size: int64Cache.Get("daily-foreclosure-refresh-input-size"),
				}),
				FileType: DataFileTypeDailyForeclosure,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeDailyForeclosurePriorityGroup,
					DataFileTypeDailyForeclosureRefreshPriorityGroup,
					2,
					DataFileTypeDailyForeclosurePriority,
				},
			},
		},
		{
			name: "hoa",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("hoa-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("hoa-input-name", "RENTNECTAR_HOA_0027.zip"),
					Size: int64Cache.SetGet("hoa-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("hoa-input-dir"),
					Name: stringCache.Get("hoa-input-name"),
					Size: int64Cache.Get("hoa-input-size"),
				}),
				FileType: DataFileTypeHOA,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeHOAPriority,
					27,
				},
			},
		},
		{
			name: "listings-before-316",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("listings-before-316-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("listings-before-316-input-name", "RENTNECTAR_LISTINGANALYTICSCOMPLETE_0124.zip"),
					Size: int64Cache.SetGet("listings-before-316-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("listings-before-316-input-dir"),
					Name: stringCache.Get("listings-before-316-input-name"),
					Size: int64Cache.Get("listings-before-316-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "listings-before-347",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("listings-before-347-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("listings-before-347-input-name", "RENTNECTAR_LISTINGANALYTICSCOMPLETE_0316_001.zip"),
					Size: int64Cache.SetGet("listings-before-347-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("listings-before-347-input-dir"),
					Name: stringCache.Get("listings-before-347-input-name"),
					Size: int64Cache.Get("listings-before-347-input-size"),
				}),
				FileType: DataFileTypeListing,
				Priorities: []int32{
					DataFileTypeListingPriority,
					316,
				},
			},
		},
		{
			name: "listings",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("listings-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("listings-input-name", "RENTNECTAR_LISTINGANALYTICSCOMPLETE_0347.zip"),
					Size: int64Cache.SetGet("listings-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("listings-input-dir"),
					Name: stringCache.Get("listings-input-name"),
					Size: int64Cache.Get("listings-input-size"),
				}),
				FileType: DataFileTypeListingV20250417,
				Priorities: []int32{
					DataFileTypeListingPriority,
					347,
				},
			},
		},
		{
			name: "monthly-amortized-equity",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("monthly-amortized-equity-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("monthly-amortized-equity-input-name", "RENTNECTAR_MONTHLY_AMORTIZEDEQUITY_0019.zip"),
					Size: int64Cache.SetGet("monthly-amortized-equity-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("monthly-amortized-equity-input-dir"),
					Name: stringCache.Get("monthly-amortized-equity-input-name"),
					Size: int64Cache.Get("monthly-amortized-equity-input-size"),
				}),
				FileType: DataFileTypeMonthlyAmortizedEquity,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeMonthlyAmortizedEquityPriority,
					19,
				},
			},
		},
		{
			name: "monthly-loan-model",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("monthly-loan-model-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("monthly-loan-model-input-name", "RENTNECTAR_MONTHLY_LOANMODEL_0020.zip"),
					Size: int64Cache.SetGet("monthly-loan-model-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("monthly-loan-model-input-dir"),
					Name: stringCache.Get("monthly-loan-model-input-name"),
					Size: int64Cache.Get("monthly-loan-model-input-size"),
				}),
				FileType: DataFileTypeMonthlyLoanModel,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeMonthlyLoanModelPriority,
					20,
				},
			},
		},
		{
			name: "property-deletes-before-242",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("property-deletes-before-242-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("property-deletes-before-242-input-name", "RENTNECTAR_PROPERTYDELETES_0091.zip"),
					Size: int64Cache.SetGet("property-deletes-before-242-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("property-deletes-before-242-input-dir"),
					Name: stringCache.Get("property-deletes-before-242-input-name"),
					Size: int64Cache.Get("property-deletes-before-242-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "property-deletes",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("property-deletes-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("property-deletes-input-name", "RENTNECTAR_PROPERTYDELETES_0242.zip"),
					Size: int64Cache.SetGet("property-deletes-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("property-deletes-input-dir"),
					Name: stringCache.Get("property-deletes-input-name"),
					Size: int64Cache.Get("property-deletes-input-size"),
				}),
				FileType: DataFileTypePropertyDelete,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					DataFileTypeAssessorRegularPriorityGroup,
					242,
					DataFileTypePropertyDeletesRegularPriority,
				},
			},
		},
		{
			name: "property-deletes-refresh",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("property-deletes-refresh-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("property-deletes-refresh-input-name", "RENTNECTAR_REFRESH_PROPERTYDELETES_0001_001.zip"),
					Size: int64Cache.SetGet("property-deletes-refresh-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("property-deletes-refresh-input-dir"),
					Name: stringCache.Get("property-deletes-refresh-input-name"),
					Size: int64Cache.Get("property-deletes-refresh-input-size"),
				}),
				FileType: DataFileTypePropertyDelete,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					DataFileTypeAssessorRefreshPriorityGroup,
					1,
					DataFileTypePropertyDeletesRefreshPriority,
				},
			},
		},
		{
			name: "recorder-before-246",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-before-246-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-before-246-input-name", "RENTNECTAR_RECORDER_0093.zip"),
					Size: int64Cache.SetGet("recorder-before-246-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-before-246-input-dir"),
					Name: stringCache.Get("recorder-before-246-input-name"),
					Size: int64Cache.Get("recorder-before-246-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "recorder",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-input-name", "RENTNECTAR_RECORDER_0246.zip"),
					Size: int64Cache.SetGet("recorder-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-input-dir"),
					Name: stringCache.Get("recorder-input-name"),
					Size: int64Cache.Get("recorder-input-size"),
				}),
				FileType: DataFileTypeRecorder,
				Priorities: []int32{
					DataFileTypeRecorderPriorityGroup,
					DataFileTypeRecorderRegularPriorityGroup,
					246,
					DataFileTypeRecorderRegularPriority,
				},
			},
		},
		{
			name: "recorder-refresh-ignore-1",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-refresh-ignore-1-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-refresh-ignore-1-input-name", "RENTNECTAR_REFRESH_RECORDER_0001_001.zip"),
					Size: int64Cache.SetGet("recorder-refresh-ignore-1-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-refresh-ignore-1-input-dir"),
					Name: stringCache.Get("recorder-refresh-ignore-1-input-name"),
					Size: int64Cache.Get("recorder-refresh-ignore-1-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "recorder-refresh",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-refresh-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-refresh-input-name", "RENTNECTAR_REFRESH_RECORDER_0002_001.zip"),
					Size: int64Cache.SetGet("recorder-refresh-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-refresh-input-dir"),
					Name: stringCache.Get("recorder-refresh-input-name"),
					Size: int64Cache.Get("recorder-refresh-input-size"),
				}),
				FileType: DataFileTypeRecorder,
				Priorities: []int32{
					DataFileTypeRecorderPriorityGroup,
					DataFileTypeRecorderRefreshPriorityGroup,
					2,
					DataFileTypeRecorderRefreshPriority,
				},
			},
		},
		{
			name: "recorder-deletes-before-246",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-deletes-before-246-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-deletes-before-246-input-name", "RENTNECTAR_RECORDERDELETES_0095.zip"),
					Size: int64Cache.SetGet("recorder-deletes-before-246-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-deletes-before-246-input-dir"),
					Name: stringCache.Get("recorder-deletes-before-246-input-name"),
					Size: int64Cache.Get("recorder-deletes-before-246-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "recorder-deletes",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-deletes-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-deletes-input-name", "RENTNECTAR_RECORDERDELETES_0246.zip"),
					Size: int64Cache.SetGet("recorder-deletes-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-deletes-input-dir"),
					Name: stringCache.Get("recorder-deletes-input-name"),
					Size: int64Cache.Get("recorder-deletes-input-size"),
				}),
				FileType: DataFileTypeRecorderDelete,
				Priorities: []int32{
					DataFileTypeRecorderPriorityGroup,
					DataFileTypeRecorderRegularPriorityGroup,
					246,
					DataFileTypeRecorderDeletesRegularPriority,
				},
			},
		},
		{
			name: "recorder-deletes-refresh-ignore-1",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-deletes-refresh-ignore-1-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-deletes-refresh-ignore-1-input-name", "RENTNECTAR_REFRESH_RECORDERDELETES_0001_001.zip"),
					Size: int64Cache.SetGet("recorder-deletes-refresh-ignore-1-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-deletes-refresh-ignore-1-input-dir"),
					Name: stringCache.Get("recorder-deletes-refresh-ignore-1-input-name"),
					Size: int64Cache.Get("recorder-deletes-refresh-ignore-1-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "recorder-deletes-refresh",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("recorder-deletes-refresh-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("recorder-deletes-refresh-input-name", "RENTNECTAR_REFRESH_RECORDERDELETES_0002_001.zip"),
					Size: int64Cache.SetGet("recorder-deletes-refresh-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("recorder-deletes-refresh-input-dir"),
					Name: stringCache.Get("recorder-deletes-refresh-input-name"),
					Size: int64Cache.Get("recorder-deletes-refresh-input-size"),
				}),
				FileType: DataFileTypeRecorderDelete,
				Priorities: []int32{
					DataFileTypeRecorderPriorityGroup,
					DataFileTypeRecorderRefreshPriorityGroup,
					2,
					DataFileTypeRecorderDeletesRefreshPriority,
				},
			},
		},
		{
			name: "rental-avm-before-40",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("rental-avm-before-40-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("rental-avm-before-40-input-name", "RENTNECTAR_RENTALAVM_0019.zip"),
					Size: int64Cache.SetGet("rental-avm-before-40-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("rental-avm-before-40-input-dir"),
					Name: stringCache.Get("rental-avm-before-40-input-name"),
					Size: int64Cache.Get("rental-avm-before-40-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "rental-avm",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("rental-avm-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("rental-avm-input-name", "RENTNECTAR_RENTALAVM_0041.zip"),
					Size: int64Cache.SetGet("rental-avm-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("rental-avm-input-dir"),
					Name: stringCache.Get("rental-avm-input-name"),
					Size: int64Cache.Get("rental-avm-input-size"),
				}),
				FileType: DataFileTypeRentalAvm,
				Priorities: []int32{
					DataFileTypeRentalAvmPriority,
					41,
				},
			},
		},
		{
			name: "tax-assessor-before-242",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("tax-assessor-before-242-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("tax-assessor-before-242-input-name", "RENTNECTAR_TAXASSESSOR_0096.zip"),
					Size: int64Cache.SetGet("tax-assessor-before-242-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("tax-assessor-before-242-input-dir"),
					Name: stringCache.Get("tax-assessor-before-242-input-name"),
					Size: int64Cache.Get("tax-assessor-before-242-input-size"),
				}),
				Ignore: true,
			},
		},
		{
			name: "tax-assessor",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("tax-assessor-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("tax-assessor-input-name", "RENTNECTAR_TAXASSESSOR_0242.zip"),
					Size: int64Cache.SetGet("tax-assessor-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("tax-assessor-input-dir"),
					Name: stringCache.Get("tax-assessor-input-name"),
					Size: int64Cache.Get("tax-assessor-input-size"),
				}),
				FileType: DataFileTypeAssessor,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					DataFileTypeAssessorRegularPriorityGroup,
					242,
					DataFileTypeAssessorRegularPriority,
				},
			},
		},
		{
			name: "tax-assessor-refresh",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("tax-assessor-refresh-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("tax-assessor-refresh-input-name", "RENTNECTAR_REFRESH_TAXASSESSOR_0001_001.zip"),
					Size: int64Cache.SetGet("tax-assessor-refresh-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("tax-assessor-refresh-input-dir"),
					Name: stringCache.Get("tax-assessor-refresh-input-name"),
					Size: int64Cache.Get("tax-assessor-refresh-input-size"),
				}),
				FileType: DataFileTypeAssessor,
				Priorities: []int32{
					DataFileTypeAssessorPriorityGroup,
					DataFileTypeAssessorRefreshPriorityGroup,
					1,
					DataFileTypeAssessorRefreshPriority,
				},
			},
		},
		{
			name: "xref-property-boundary-match-parcel",
			in: &entities.CreateDataFileEntryInput{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.SetGet("xref-property-boundary-match-parcel-input-dir", uuid.New().String()),
					Name: stringCache.SetGet("xref-property-boundary-match-parcel-input-name", "RENTNECTAR_XREF_PROPERTYTOBOUNDARYMATCH_PARCEL_0004.zip"),
					Size: int64Cache.SetGet("xref-property-boundary-match-parcel-input-size", must.GenRandomInt64()),
				}),
			},
			out: &entities.DataFileEntry{
				StorageObject: storage.NewObject(&storage.NewObjectInput{
					Dir:  stringCache.Get("xref-property-boundary-match-parcel-input-dir"),
					Name: stringCache.Get("xref-property-boundary-match-parcel-input-name"),
					Size: int64Cache.Get("xref-property-boundary-match-parcel-input-size"),
				}),
				FileType: DataFileTypeXrefPropertyToBoundaryMatchParcel,
				Ignore:   true,
				Priorities: []int32{
					DataFileTypeXrefPropertyToBoundaryMatchParcelPriority,
					4,
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
