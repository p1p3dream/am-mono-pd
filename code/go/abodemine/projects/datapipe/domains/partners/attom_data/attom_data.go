package attom_data

import (
	"path"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/projects/datapipe/entities"
)

var PartnerId = uuid.MustParse("11ecadd9-6bc1-4b5d-927e-378dc11829fb")

const (
	// Use random 9-digit integers (32bit) to ensure new
	// data file types can be added and grouped together
	// with similar types without worrying about
	// collisions and order.
	//
	// Although smaller values would be easier to remember,
	// larger values are used to reduce the chance of
	// collisions with data file types from other partners.
	//
	// New values can be generated with: ugen --digit -l 9.
	DataFileTypeAssessor                  entities.DataFileType = 369011232
	DataFileTypeAssignmentRelease         entities.DataFileType = 466086029
	DataFileTypeAvm                       entities.DataFileType = 425540929
	DataFileTypeBuildingPermitClassifiers entities.DataFileType = 554513072
	DataFileTypeBuildingPermitDelete      entities.DataFileType = 908335649
	DataFileTypeBuildingPermitStatus      entities.DataFileType = 902471889
	DataFileTypeBuildingPermit            entities.DataFileType = 210907693
	DataFileTypeCF                        entities.DataFileType = 929659292
	DataFileTypeDailyForeclosure          entities.DataFileType = 779874070
	DataFileTypeHOA                       entities.DataFileType = 313638910

	// Ensure future versions of the same file type ALWAYS
	// have a higher value than the previous one so that
	// we can easily order their processing.
	DataFileTypeListing entities.DataFileType = 233079812
	// Added the CurrentStatus (YN) field.
	DataFileTypeListingV20250417 entities.DataFileType = 233079813

	DataFileTypeMonthlyAmortizedEquity            entities.DataFileType = 702256699
	DataFileTypeMonthlyLoanModel                  entities.DataFileType = 969505053
	DataFileTypePropertyDelete                    entities.DataFileType = 712329206
	DataFileTypeRecorder                          entities.DataFileType = 260494762
	DataFileTypeRecorderDelete                    entities.DataFileType = 990733446
	DataFileTypeRentalAvm                         entities.DataFileType = 614553397
	DataFileTypeXrefPropertyToBoundaryMatchParcel entities.DataFileType = 955541030
)

const (
	DataFileTypeListingPriority int32 = 1 + iota
	DataFileTypeListingV20250417Priority

	// Use a priority group for recorder files to ensure
	// that deletes are processed before regular updates.
	// This is reversed for refresh files.
	DataFileTypeRecorderPriorityGroup
	DataFileTypeRecorderRefreshPriorityGroup
	DataFileTypeRecorderRegularPriorityGroup
	DataFileTypeRecorderRefreshPriority
	DataFileTypeRecorderDeletesRefreshPriority
	DataFileTypeRecorderDeletesRegularPriority
	DataFileTypeRecorderRegularPriority

	// Use a priority group for assessor files to ensure
	// that deletes are processed before regular updates.
	// This is reversed for refresh files.
	DataFileTypeAssessorPriorityGroup
	DataFileTypeAssessorRefreshPriorityGroup
	DataFileTypeAssessorRegularPriorityGroup
	DataFileTypeAssessorRefreshPriority
	DataFileTypePropertyDeletesRefreshPriority
	DataFileTypePropertyDeletesRegularPriority
	DataFileTypeAssessorRegularPriority

	DataFileTypeAssignmentReleasePriority
	DataFileTypeAvmPriority
	DataFileTypeBuildingPermitClassifiersPriority
	DataFileTypeBuildingPermitDeletePriority
	DataFileTypeBuildingPermitStatusPriority
	DataFileTypeBuildingPermitPriority
	DataFileTypeCFPriority

	DataFileTypeDailyForeclosurePriorityGroup
	DataFileTypeDailyForeclosureRefreshPriorityGroup
	DataFileTypeDailyForeclosureRegularPriorityGroup
	DataFileTypeDailyForeclosurePriority

	DataFileTypeHOAPriority
	DataFileTypeMonthlyAmortizedEquityPriority
	DataFileTypeMonthlyLoanModelPriority
	DataFileTypeRentalAvmPriority
	DataFileTypeXrefPropertyToBoundaryMatchParcelPriority
)

type DataSource struct{}

func (ds *DataSource) CreateDataFileEntry(r *arc.Request, in *entities.CreateDataFileEntryInput) (*entities.DataFileEntry, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "5e376d72-7d2d-4a0b-a731-421d45ba5b75",
			Code:   errors.Code_INTERNAL,
			Detail: "Input is nil.",
		}
	}

	obj := in.StorageObject

	if obj == nil {
		return nil, &errors.Object{
			Id:     "2f63170e-74a3-4cb6-b5ba-7f7519138634",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing storage object.",
		}
	}

	out := &entities.DataFileEntry{
		StorageObject: in.StorageObject,
	}

	if obj.IsDirectory() {
		switch {
		case path.Dir(in.Path) == "/ftp" && obj.Name == "Refresh":
			// Directory for attom data refresh files, created 2025-03.

			out.FileType = entities.DataFileTypeSelectedDirectory
			out.EnterDirectory = true
			out.IgnoreSubDirs = true
		default:
			// Skip other directories since Attom Data just put files in the root.
			out.Ignore = true
		}

		return out, nil
	}

	standardMatch := regexp.MustCompile(`^([A-Za-z]+)_([A-Za-z_]+)_(\d+)_?(\d+)?.zip$`)

	matches := standardMatch.FindStringSubmatch(obj.Name)
	if len(matches) == 0 {
		log.Debug().
			Str("name", obj.Name).
			Str("reason", "no_match").
			Msg("Skipping file.")

		out.Ignore = true
		return out, nil
	}

	if len(matches) < 4 {
		log.Warn().
			Str("name", obj.Name).
			Str("reason", "missing_submatches").
			Msg("Skipping file.")

		out.Ignore = true
		return out, nil
	}

	releaseNumber, err := strconv.ParseInt(matches[3], 10, 32)
	if err != nil {
		log.Warn().
			Str("name", obj.Name).
			Str("reason", "invalid_release_number").
			Str("release_number", matches[3]).
			Msg("Skipping file.")

		out.Ignore = true
		return out, nil
	}
	out.ReleaseNumber = int32(releaseNumber)

	if len(matches) > 4 && matches[4] != "" {
		// NOTICE: out.ReleasePart is irrelevant for out.Priorities.
		releasePart, err := strconv.ParseInt(matches[4], 10, 32)
		if err != nil {
			log.Warn().
				Str("name", obj.Name).
				Str("reason", "invalid_release_part").
				Int32("release_number", out.ReleaseNumber).
				Str("release_part", matches[4]).
				Msg("Skipping file.")

			out.Ignore = true
			return out, nil
		}
		out.ReleasePart = int32(releasePart)
	}

	fileTypeName := matches[2]
	isRefreshFile := path.Dir(in.Path) == "/ftp/Refresh"

	switch fileTypeName {
	case "ASSIGNMENTRELEASE":
		out.Ignore = true
		out.FileType = DataFileTypeAssignmentRelease
		out.Priorities = []int32{
			DataFileTypeAssignmentReleasePriority,
			out.ReleaseNumber,
		}
	case "AVM":
		out.Ignore = true
		out.FileType = DataFileTypeAvm
		out.Priorities = []int32{
			DataFileTypeAvmPriority,
			out.ReleaseNumber,
		}
	case "BUILDINGPERMIT":
		out.Ignore = true
		out.FileType = DataFileTypeBuildingPermit
		out.Priorities = []int32{
			DataFileTypeBuildingPermitPriority,
			out.ReleaseNumber,
		}
	case "BUILDINGPERMITCLASSIFIERS":
		out.Ignore = true
		out.FileType = DataFileTypeBuildingPermitClassifiers
		out.Priorities = []int32{
			DataFileTypeBuildingPermitClassifiersPriority,
			out.ReleaseNumber,
		}
	case "BUILDINGPERMITDELETE":
		out.Ignore = true
		out.FileType = DataFileTypeBuildingPermitDelete
		out.Priorities = []int32{
			DataFileTypeBuildingPermitDeletePriority,
			out.ReleaseNumber,
		}
	case "BUILDINGPERMITSTATUS":
		out.Ignore = true
		out.FileType = DataFileTypeBuildingPermitStatus
		out.Priorities = []int32{
			DataFileTypeBuildingPermitStatusPriority,
			out.ReleaseNumber,
		}
	case "CF":
		out.Ignore = true
		out.FileType = DataFileTypeCF
		out.Priorities = []int32{
			DataFileTypeCFPriority,
			out.ReleaseNumber,
		}
	case "DAILY_FORECLOSURE":
		out.Ignore = true
		out.FileType = DataFileTypeDailyForeclosure
		out.Priorities = []int32{
			DataFileTypeDailyForeclosurePriorityGroup,
			DataFileTypeDailyForeclosureRegularPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeDailyForeclosurePriority,
		}
	case "REFRESH_FORECLOSURE":
		out.Ignore = true
		out.FileType = DataFileTypeDailyForeclosure
		out.Priorities = []int32{
			DataFileTypeDailyForeclosurePriorityGroup,
			DataFileTypeDailyForeclosureRefreshPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeDailyForeclosurePriority,
		}
	case "HOA":
		out.Ignore = true
		out.FileType = DataFileTypeHOA
		out.Priorities = []int32{
			DataFileTypeHOAPriority,
			out.ReleaseNumber,
		}
	case "LISTINGANALYTICSCOMPLETE":
		if out.ReleaseNumber < 316 {
			out.Ignore = true
			return out, nil
		}

		if out.ReleaseNumber < 347 {
			out.FileType = DataFileTypeListing
		} else {
			out.FileType = DataFileTypeListingV20250417
		}

		out.Priorities = []int32{
			DataFileTypeListingPriority,
			out.ReleaseNumber,
		}
	case "MONTHLY_AMORTIZEDEQUITY":
		out.Ignore = true
		out.FileType = DataFileTypeMonthlyAmortizedEquity
		out.Priorities = []int32{
			DataFileTypeMonthlyAmortizedEquityPriority,
			out.ReleaseNumber,
		}
	case "MONTHLY_LOANMODEL":
		out.Ignore = true
		out.FileType = DataFileTypeMonthlyLoanModel
		out.Priorities = []int32{
			DataFileTypeMonthlyLoanModelPriority,
			out.ReleaseNumber,
		}
	case "PROPERTYDELETES":
		if out.ReleaseNumber < 242 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypePropertyDelete
		out.Priorities = []int32{
			DataFileTypeAssessorPriorityGroup,
			DataFileTypeAssessorRegularPriorityGroup,
			out.ReleaseNumber,
			DataFileTypePropertyDeletesRegularPriority,
		}
	case "REFRESH_PROPERTYDELETES":
		out.FileType = DataFileTypePropertyDelete
		out.Priorities = []int32{
			DataFileTypeAssessorPriorityGroup,
			DataFileTypeAssessorRefreshPriorityGroup,
			out.ReleaseNumber,
			DataFileTypePropertyDeletesRefreshPriority,
		}
	case "RECORDER":
		if out.ReleaseNumber < 246 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeRecorder
		out.Priorities = []int32{
			DataFileTypeRecorderPriorityGroup,
			DataFileTypeRecorderRegularPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeRecorderRegularPriority,
		}
	case "REFRESH_RECORDER":
		if out.ReleaseNumber != 2 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeRecorder
		out.Priorities = []int32{
			DataFileTypeRecorderPriorityGroup,
			DataFileTypeRecorderRefreshPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeRecorderRefreshPriority,
		}
	case "RECORDERDELETES":
		if out.ReleaseNumber < 246 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeRecorderDelete
		out.Priorities = []int32{
			DataFileTypeRecorderPriorityGroup,
			DataFileTypeRecorderRegularPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeRecorderDeletesRegularPriority,
		}
	case "REFRESH_RECORDERDELETES":
		if out.ReleaseNumber != 2 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeRecorderDelete
		out.Priorities = []int32{
			DataFileTypeRecorderPriorityGroup,
			DataFileTypeRecorderRefreshPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeRecorderDeletesRefreshPriority,
		}
	case "RENTALAVM":
		if !isRefreshFile && out.ReleaseNumber < 40 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeRentalAvm
		out.Priorities = []int32{
			DataFileTypeRentalAvmPriority,
			out.ReleaseNumber,
		}
	case "REFRESH_RENTALAVM":
		out.Ignore = true
		return out, nil
	case "TAXASSESSOR":
		if out.ReleaseNumber < 242 {
			out.Ignore = true
			return out, nil
		}

		out.FileType = DataFileTypeAssessor
		out.Priorities = []int32{
			DataFileTypeAssessorPriorityGroup,
			DataFileTypeAssessorRegularPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeAssessorRegularPriority,
		}
	case "REFRESH_TAXASSESSOR":
		out.FileType = DataFileTypeAssessor
		out.Priorities = []int32{
			DataFileTypeAssessorPriorityGroup,
			DataFileTypeAssessorRefreshPriorityGroup,
			out.ReleaseNumber,
			DataFileTypeAssessorRefreshPriority,
		}
	case "XREF_PROPERTYTOBOUNDARYMATCH_PARCEL":
		out.Ignore = true
		out.FileType = DataFileTypeXrefPropertyToBoundaryMatchParcel
		out.Priorities = []int32{
			DataFileTypeXrefPropertyToBoundaryMatchParcelPriority,
			out.ReleaseNumber,
		}
	default:
		log.Warn().
			Str("filename", obj.Name).
			Msg("Unknown data file type.")
	}

	return out, nil
}

func (ds *DataSource) DataRecordByFileType(fileType entities.DataFileType) (entities.DataRecord, error) {
	switch fileType {
	case DataFileTypeAssessor:
		return new(Assessor), nil
	case DataFileTypeListing:
		return new(Listing), nil
	case DataFileTypeListingV20250417:
		return &Listing{
			fileType: fileType,
		}, nil
	case DataFileTypePropertyDelete:
		return new(PropertyDelete), nil
	case DataFileTypeRecorder:
		return new(Recorder), nil
	case DataFileTypeRecorderDelete:
		return new(RecorderDelete), nil
	case DataFileTypeRentalAvm:
		return new(RentalAvm), nil
	}

	return nil, &errors.Object{
		Id:     "7c543b36-018b-4d1e-a3c8-260d89e3e4ff",
		Code:   errors.Code_INTERNAL,
		Detail: "Unknown data file type.",
		Meta: map[string]any{
			"file_type": fileType,
		},
	}
}

func (ds *DataSource) FieldSeparatorByFileType(fileType entities.DataFileType) string {
	return "\t"
}
