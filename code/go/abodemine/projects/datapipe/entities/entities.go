package entities

import (
	"bufio"
	"time"

	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/lib/storage"
)

type DataFileType int32

const (
	// Represents a directory that was specifically
	// selected for processing by the DataSource.
	DataFileTypeSelectedDirectory DataFileType = 292376608
)

const (
	DataFileDirectoryStatusToDo       = 100
	DataFileDirectoryStatusInProgress = 200
	DataFileDirectoryStatusDone       = 300
	DataFileDirectoryStatusIgnored    = 400
)

type DataFileDirectory struct {
	Id        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Meta      map[string]any `json:"meta"`

	PartnerId         uuid.UUID  `json:"partner_id"`
	ParentDirectoryId *uuid.UUID `json:"parent_directory_id"`
	Status            int32      `json:"status"`
	Path              string     `json:"path"`
	Name              string     `json:"name"`

	// Priorities is a list that is used to determine
	// the order of processing, from highest to lowest.
	Priorities []int32 `json:"priorities"`
}

const (
	DataFileObjectStatusToDo       = 100
	DataFileObjectStatusInProgress = 200
	DataFileObjectStatusDone       = 300
	DataFileObjectStatusIgnored    = 400
)

type DataRecord interface {
	New(headers map[int]string, fields []string) (DataRecord, error)
	LoadParams() *DataRecordLoadParams
	SQLColumns() []string
	SQLTable() string
	SQLValues() ([]any, error)
}

const (
	// Ensure that mode must be explicitly set.
	DataRecordModeUndefined = iota

	// For DataFile types that can be deleted in bulk.
	// E.g. RecorderDelete.
	DataRecordModeBatchDelete

	// For DataFile types that can be inserted in bulk.
	// E.g. AVMs, Recorder.
	DataRecordModeBatchInsert

	// For DataFile types that require a custom load function.
	// E.g. Assessor, Listing.
	DataRecordModeLoadFunc
)

type DataRecordLoadParams struct {
	LoadFunc func(r *arc.Request, in *LoadDataRecordInput) (*LoadDataRecordOutput, error)
	Mode     int
}

type LoadDataRecordInput struct {
	BatchSize                int32
	Columns                  []string
	DataFileObject           *DataFileObject
	DataRecord               DataRecord
	FieldSeparator           string
	Headers                  map[int]string
	Scanner                  *bufio.Scanner
	UpdateDataFileObjectFunc func(r *arc.Request, in *UpdateDataFileObjectInput) (*UpdateDataFileObjectOutput, error)
}

type LoadDataRecordOutput struct {
	DeletedRecords   int64
	ProcessedRecords int64
}

type DataFileObject struct {
	Id        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Meta      map[string]any `json:"meta"`

	DirectoryId  *uuid.UUID   `json:"directory_id"`
	ParentFileId *uuid.UUID   `json:"parent_file_id"`
	FileType     DataFileType `json:"file_type"`
	Hash         []byte       `json:"hash"`
	Status       int32        `json:"status"`
	RecordCount  int32        `json:"record_count"`

	// FileDir is the original directory name where
	// the file was located.
	FileDir string `json:"file_dir"`

	// FileName is the original file name.
	FileName string `json:"file_name"`

	// FileSize is the size of the file in bytes.
	FileSize *int64 `json:"file_size"`

	// Priorities is a list that is used to determine
	// the order of processing, from highest to lowest.
	Priorities []int32 `json:"priorities"`

	WorkerId *uuid.UUID `json:"worker_id"`
}

// DataFileEntry contains information needed to
// process DataFile directories and objects.
type DataFileEntry struct {
	FileType      DataFileType
	StorageObject *storage.Object

	// Required by FirstAmerican.
	Date time.Time

	// Required by AttomData.
	ReleaseNumber int32
	ReleasePart   int32

	// Do not process this entry.
	Ignore bool
	// Do not recurse into subdirectories.
	IgnoreSubDirs bool

	// Enter the directory even if its parent
	// is configured to IgnoreSubDirs.
	EnterDirectory bool

	Priorities []int32
}

type CreateDataFileEntryInput struct {
	Path          string
	StorageObject *storage.Object
}

type SkipEntryInput struct {
	FileType DataFileType
	Index    int
	Length   int
}

type DataSource interface {
	CreateDataFileEntry(r *arc.Request, in *CreateDataFileEntryInput) (*DataFileEntry, error)

	DataRecordByFileType(fileTYpe DataFileType) (DataRecord, error)
	FieldSeparatorByFileType(fileType DataFileType) string
}

type UpdateDataFileObjectInput struct {
	Id          uuid.UUID
	UpdatedAt   time.Time
	Meta        map[string]any
	Status      int32
	RecordCount int32
	Priorities  []int32
}

type UpdateDataFileObjectOutput struct {
	Entity *DataFileObject
}
