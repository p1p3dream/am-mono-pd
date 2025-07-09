package first_american

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/entities"
)

var PartnerId = uuid.MustParse("44f2033f-f93a-4cad-bcaa-d5f649940094")

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
	DataFileTypeASR            entities.DataFileType = 494501913
	DataFileTypeAddress        entities.DataFileType = 732748435
	DataFileTypeAssessorAnnual entities.DataFileType = 363785734
	DataFileTypeAssessorUpdate entities.DataFileType = 515042169
	DataFileTypeAVMPower       entities.DataFileType = 288491261
	DataFileTypeDeedMtg        entities.DataFileType = 541242299
	DataFileTypeHOA            entities.DataFileType = 347768023
	DataFileTypeHOALien        entities.DataFileType = 156325817
	DataFileTypeHPI            entities.DataFileType = 579722379
	DataFileTypeInvLien        entities.DataFileType = 357089610
	DataFileTypeListing        entities.DataFileType = 960605711
	DataFileTypeNOD            entities.DataFileType = 311010286
	DataFileTypePFC            entities.DataFileType = 448481785
	DataFileTypeShape          entities.DataFileType = 236998122
	DataFileTypeTaxHistory     entities.DataFileType = 240431550
	DataFileTypeValueHistory   entities.DataFileType = 809597855
)

const (
	// DataFileType by descending priority.
	// Lower values are processed first.

	// Use a priority group for assessor files to ensure
	// that annual and update are processed by date first,
	// then priority.
	DataFileTypeAssessorPriorityGroup int32 = 1 + iota
	DataFileTypeAssessorUpdatePriority
	DataFileTypeAssessorAnnualPriority

	DataFileTypeAVMPowerPriority
	DataFileTypeAVMPowerHistoryPriority
	DataFileTypeASRPriority
	DataFileTypeDeedMtgPriority
	DataFileTypeHOAPriority
	DataFileTypeHOALienPriority
	DataFileTypeHPIPriority
	DataFileTypeInvLienPriority
	DataFileTypeListingPriority
	DataFileTypeNODPriority
	DataFileTypePFCPriority
	DataFileTypeShapePriority
	DataFileTypeTaxHistoryPriority
	DataFileTypeValueHistoryPriority
	DataFileTypeAddressPriority
)

type DataSource struct{}

func (ds *DataSource) CreateDataFileEntry(r *arc.Request, in *entities.CreateDataFileEntryInput) (*entities.DataFileEntry, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "a84a2463-5000-4ec7-9646-1fbec9a93d95",
			Code:   errors.Code_INTERNAL,
			Detail: "Input is nil.",
		}
	}

	obj := in.StorageObject

	if obj == nil {
		return nil, &errors.Object{
			Id:     "c5ad6b0e-bfee-4617-8fae-672c0a93e0a2",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing storage object.",
		}
	}

	out := &entities.DataFileEntry{
		StorageObject: in.StorageObject,
	}

	// Skip files, since First American groups by directory.
	if !obj.IsDirectory() {
		return out, nil
	}

	standardMatch := regexp.MustCompile(`^([0-9]{8})_([A-Za-z_]*[A-Za-z]+)_?([0-9]+)?$`)

	matches := standardMatch.FindStringSubmatch(obj.Name)
	if len(matches) == 0 {
		log.Debug().
			Str("name", obj.Name).
			Str("reason", "no_match").
			Msg("Skipping file.")

		out.Ignore = true
		return out, nil
	}

	if len(matches) < 3 {
		log.Warn().
			Str("name", obj.Name).
			Str("reason", "missing_submatches").
			Msg("Skipping file.")

		out.Ignore = true
		return out, nil
	}

	var err error

	out.Date, err = time.Parse(consts.IntegerDate, matches[1])
	if err != nil {
		return nil, &errors.Object{
			Id:     "0f889387-e353-4bb4-b2f0-e93a80191c3e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid date format.",
			Meta: map[string]any{
				"date": matches[1],
			},
		}
	}

	if len(matches) > 3 && matches[3] != "" {
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
	}

	iDate := int32(val.IntegerDate(out.Date))
	suffix := strings.ToUpper(matches[2])

	switch suffix {
	case "ADDRESS", "ADDRESSMASTER":
		// We won't process master address until we are sure there
		// won't be conflicts with our own addresses table.

		out.Ignore = true
		out.FileType = DataFileTypeAddress

		out.Priorities = []int32{
			DataFileTypeAddressPriority,
			iDate,
		}
	case "ANNUAL":
		out.FileType = DataFileTypeAssessorAnnual

		if iDate < 20250326 {
			out.Ignore = true
		} else {
			out.Priorities = []int32{
				DataFileTypeAssessorPriorityGroup,
				iDate,
				DataFileTypeAssessorAnnualPriority,
			}
		}
	case "ASR":
		out.Ignore = true
		out.FileType = DataFileTypeASR

		out.Priorities = []int32{
			DataFileTypeASRPriority,
			iDate,
		}
	case "AVMPOWER":
		if out.Date.Before(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) {
			out.Ignore = true
		} else {
			out.FileType = DataFileTypeAVMPower
			out.Priorities = []int32{
				DataFileTypeAVMPowerPriority,
				iDate,
			}
		}
	case "DEED", "DEEDMTG":
		out.Ignore = true
		out.FileType = DataFileTypeDeedMtg

		out.Priorities = []int32{
			DataFileTypeDeedMtgPriority,
			iDate,
		}
	case "HOA":
		out.Ignore = true
		out.FileType = DataFileTypeHOA

		out.Priorities = []int32{
			DataFileTypeHOAPriority,
			iDate,
		}
	case "HOALIEN":
		out.Ignore = true
		out.FileType = DataFileTypeHOALien

		out.Priorities = []int32{
			DataFileTypeHOALienPriority,
			iDate,
		}
	case "HPI":
		out.Ignore = true
		out.FileType = DataFileTypeHPI

		out.Priorities = []int32{
			DataFileTypeHPIPriority,
			iDate,
		}
	case "INVL", "INVLIEN":
		out.Ignore = true
		out.FileType = DataFileTypeInvLien

		out.Priorities = []int32{
			DataFileTypeInvLienPriority,
			iDate,
		}
	case "LISTING", "LISTINGS":
		out.Ignore = true
		out.FileType = DataFileTypeListing

		out.Priorities = []int32{
			DataFileTypeListingPriority,
			iDate,
		}
	case "NOD":
		out.Ignore = true
		out.FileType = DataFileTypeNOD

		out.Priorities = []int32{
			DataFileTypeNODPriority,
			iDate,
		}
	case "PFC":
		out.Ignore = true
		out.FileType = DataFileTypePFC

		out.Priorities = []int32{
			DataFileTypePFCPriority,
			iDate,
		}
	case "POWER_AVMHIST":
		// TODO.
		out.Ignore = true
		out.FileType = DataFileTypeAVMPower
		out.Priorities = []int32{
			DataFileTypeAVMPowerHistoryPriority,
			out.ReleaseNumber,
			iDate,
		}
	case "PROP":
		// These directories were (manually) created from files from either
		// Annual or Update, and should be ignored.
		out.Ignore = true
	case "SHAPE":
		out.Ignore = true
		out.FileType = DataFileTypeShape

		out.Priorities = []int32{
			DataFileTypeShapePriority,
			iDate,
		}
	case "TAXHIST", "TAXHISTORY":
		out.Ignore = true
		out.FileType = DataFileTypeTaxHistory

		out.Priorities = []int32{
			DataFileTypeTaxHistoryPriority,
			iDate,
		}
	case "UPDATE":
		out.FileType = DataFileTypeAssessorUpdate

		if iDate < 20250403 {
			out.Ignore = true
		} else {
			out.Priorities = []int32{
				DataFileTypeAssessorPriorityGroup,
				iDate,
				DataFileTypeAssessorUpdatePriority,
			}
		}
	case "VALHIST", "VALUEHIST":
		out.Ignore = true
		out.FileType = DataFileTypeValueHistory

		out.Priorities = []int32{
			DataFileTypeValueHistoryPriority,
			iDate,
		}
	default:
		log.Warn().
			Str("filename", obj.Name).
			Str("suffix", suffix).
			Msg("Unknown data file type.")
	}

	return out, nil
}

func (ds *DataSource) DataRecordByFileType(fileType entities.DataFileType) (entities.DataRecord, error) {
	switch fileType {
	case DataFileTypeAddress:
		return new(Address), nil
	case DataFileTypeAssessorAnnual:
		return new(Assessor), nil
	case DataFileTypeAssessorUpdate:
		return new(Assessor), nil
	case DataFileTypeAVMPower:
		return new(AVMPower), nil
	}

	return nil, &errors.Object{
		Id:     "a2e5f9fa-22d8-4ef4-b29d-54ac3d8ab294",
		Code:   errors.Code_INTERNAL,
		Detail: "Unknown data file type.",
		Meta: map[string]any{
			"file_type": fileType,
		},
	}
}

func (ds *DataSource) FieldSeparatorByFileType(fileType entities.DataFileType) string {
	return "|"
}
