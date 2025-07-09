package worker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/lib/distsync"
	"abodemine/lib/errors"
	"abodemine/lib/storage"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/conf"
	"abodemine/projects/datapipe/domains/partners/attom_data"
	"abodemine/projects/datapipe/domains/partners/first_american"
	"abodemine/projects/datapipe/entities"
	"abodemine/repositories/opensearch"
)

type Domain interface {
	FetchDataSource(r *arc.Request, in *FetchDataSourceInput) (*FetchDataSourceOutput, error)
	ProcessDataSource(r *arc.Request, in *ProcessDataSourceInput) (*ProcessDataSourceOutput, error)

	InsertDataFileDirectory(r *arc.Request, in *InsertDataFileDirectoryInput) (*InsertDataFileDirectoryOutput, error)
	SelectDataFileDirectory(r *arc.Request, in *SelectDataFileDirectoryInput) (*SelectDataFileDirectoryOutput, error)
	UpdateDataFileDirectory(r *arc.Request, in *UpdateDataFileDirectoryInput) (*UpdateDataFileDirectoryOutput, error)

	EnsureDataFileObject(r *arc.Request, in *EnsureDataFileObjectInput) (*EnsureDataFileObjectOutput, error)
	InsertDataFileObject(r *arc.Request, in *InsertDataFileObjectInput) (*InsertDataFileObjectOutput, error)
	SelectDataFileObject(r *arc.Request, in *SelectDataFileObjectInput) (*SelectDataFileObjectOutput, error)
	UpdateDataFileObject(r *arc.Request, in *entities.UpdateDataFileObjectInput) (*entities.UpdateDataFileObjectOutput, error)

	LoadOpenSearch(r *arc.Request, in *LoadOpenSearchInput) (*LoadOpenSearchOutput, error)

	SyncProperties(r *arc.Request, in *SyncPropertiesInput) (*SyncPropertiesOutput, error)
}

type domain struct {
	config *conf.Config

	addressDomain address.Domain

	repository         Repository
	osSearchRepository opensearch.Repository
}

type NewDomainInput struct {
	Config *conf.Config

	AddressDomain address.Domain

	Repository         Repository
	OsSearchRepository opensearch.Repository
}

func NewDomain(in *NewDomainInput) *domain {
	return &domain{
		config:        in.Config,
		addressDomain: in.AddressDomain,
		repository: val.Ternary(
			in.Repository == nil,
			NewRepository(),
			in.Repository,
		),
		osSearchRepository: in.OsSearchRepository,
	}
}

type FetchDataSourceInput struct {
	PartnerId uuid.UUID

	RcloneCheckers    string
	RcloneTransfers   string
	RcloneDestination string
	RcloneSource      string

	NoLock bool
}

type FetchDataSourceOutput struct{}

func (dom *domain) FetchDataSource(r *arc.Request, in *FetchDataSourceInput) (*FetchDataSourceOutput, error) {
	if !in.NoLock {
		lockOut, err := dom.Lock(r, &LockInput{
			PartnerId:  in.PartnerId,
			LockerName: "fetcher",
		})
		if err != nil {
			return nil, errors.Forward(err, "e500eb47-5f78-47c3-853d-40ddd64d952e")
		}

		defer func() {
			lockOut.ExtendCancel()
			lockOut.LockerWg.Wait()
		}()
	}

	rcloneDestination := in.RcloneDestination

	if rcloneDestination == "" {
		return nil, &errors.Object{
			Id:     "8bf7d35b-a20e-43f0-9601-39e934be1279",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing rclone destination.",
		}
	}

	rcloneSource := in.RcloneSource

	if rcloneSource == "" {
		return nil, &errors.Object{
			Id:     "ea36b334-9f92-443b-9200-0c7f0cc24be7",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing rclone source.",
		}
	}

	_, err := dom.ExecFetchDataSource(r, &ExecFetchDataSourceInput{
		RcloneCheckers:    in.RcloneCheckers,
		RcloneTransfers:   in.RcloneTransfers,
		RcloneSource:      rcloneSource,
		RcloneDestination: rcloneDestination,
	})
	if err != nil {
		return nil, errors.Forward(err, "2da22e57-aa3b-49a1-b475-f6246d884417")
	}

	out := &FetchDataSourceOutput{}

	return out, nil
}

type ExecFetchDataSourceInput struct {
	RcloneCheckers    string
	RcloneTransfers   string
	RcloneDestination string
	RcloneSource      string
}

type ExecFetchDataSourceOutput struct{}

func (dom *domain) ExecFetchDataSource(r *arc.Request, in *ExecFetchDataSourceInput) (*ExecFetchDataSourceOutput, error) {
	log.Info().
		Str("rclone_checkers", in.RcloneCheckers).
		Str("rclone_transfers", in.RcloneTransfers).
		Str("rclone_source", in.RcloneSource).
		Str("rclone_destination", in.RcloneDestination).
		Msg("Executing fetch data source.")

	logFileName, err := uuid.NewRandom()
	if err != nil {
		return nil, &errors.Object{
			Id:     "1592167c-59a8-4581-b188-b7cb937252e6",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to generate log file name.",
			Cause:  err.Error(),
		}
	}

	logFile := filepath.Join(os.TempDir(), logFileName.String())

	cmd := exec.Command(
		"rclone",
		"copy",
		"--size-only",
		fmt.Sprintf("--checkers=%s", in.RcloneCheckers),
		fmt.Sprintf("--transfers=%s", in.RcloneTransfers),
		in.RcloneSource,
		in.RcloneDestination,
		"--log-file", logFile,
		"--log-level", "info",
	)

	stderr := &bytes.Buffer{}

	cmd.Stderr = stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		if file, err := os.ReadFile(logFile); err == nil {
			fmt.Printf("LOGFILE:\n%s\n", string(file))
		}

		return nil, &errors.Object{
			Id:     "a7576940-82cf-464a-bfdd-7b0fd10418a5",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to run rclone.",
			Cause:  err.Error(),
			Meta: map[string]any{
				"stderr": stderr.String(),
			},
		}
	}

	file, err := os.Open(logFile)
	if err != nil {
		return nil, &errors.Object{
			Id:     "05c86642-3f86-40a6-ba04-8d1cd07084dd",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to open log file.",
			Cause:  err.Error(),
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	hasTransferredSomething := true
	logBuf := &bytes.Buffer{}

	for scanner.Scan() {
		line := scanner.Text()

		if logBuf.Len() == 0 {
			hasTransferredSomething = !strings.Contains(line, "There was nothing to transfer")
		}

		logBuf.WriteString(line)
		logBuf.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, &errors.Object{
			Id:     "f3b0c4a1-2d5e-4b8c-9f6d-7a0e1f2f3b8c",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to read log file.",
			Cause:  err.Error(),
		}
	}

	if logBuf.Len() > 0 {
		fmt.Printf("LOGFILE:\n%s\n", logBuf.String())
	}

	if hasTransferredSomething {
		// TODO: trigger loader posting to SQS.
		log.Info().
			Msg("Data source fetch completed.")
	} else {
		log.Info().
			Msg("There was nothing to transfer.")
	}

	out := &ExecFetchDataSourceOutput{}

	return out, nil
}

type ProcessDataSourceInput struct {
	PartnerId uuid.UUID
	Backend   storage.Backend

	// Max number of compressed files to process concurrently.
	FileBufferSize int

	// Max number of records to process in each batch operation.
	DatabaseBatchSize int

	PathPrefix string

	// If set, process only this PriorityGroup.
	PriorityGroup int32

	NoLock  bool
	Version string

	WorkerId *uuid.UUID
}

type ProcessDataSourceOutput struct{}

func (dom *domain) ProcessDataSource(r *arc.Request, in *ProcessDataSourceInput) (*ProcessDataSourceOutput, error) {
	if !in.NoLock {
		lockOut, err := dom.Lock(r, &LockInput{
			PartnerId:  in.PartnerId,
			LockerName: "loader",
		})
		if err != nil {
			return nil, errors.Forward(err, "111c4beb-8ca1-4818-9545-2de80a013555")
		}

		defer func() {
			lockOut.ExtendCancel()
			lockOut.LockerWg.Wait()
		}()
	}

	processDataSourceDirInput := &ProcessDataSourceDirInput{
		Backend:        in.Backend,
		FileBufferSize: in.FileBufferSize,
		IsRootDir:      true,
		PartnerId:      in.PartnerId,
		Path:           in.PathPrefix,
		PriorityGroup:  in.PriorityGroup,
		WorkerId:       in.WorkerId,
	}

	switch in.PartnerId {
	case uuid.Nil:
		return nil, &errors.Object{
			Id:     "d8153b71-390b-449d-94b4-e018cf52c47a",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing PartnerId.",
		}
	case attom_data.PartnerId:
		processDataSourceDirInput.DataSource = &attom_data.DataSource{}
		processDataSourceDirInput.IgnoreSubDirs = true
	case first_american.PartnerId:
		processDataSourceDirInput.DataSource = &first_american.DataSource{}
	default:
		return nil, &errors.Object{
			Id:     "cd37d2d5-ab8b-43e5-8e05-35ef8a20021d",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Unknown PartnerId.",
		}
	}

	switch in.Version {
	case "v1":
		// v1 was replaced by v2 and is not supported anymore.
		return nil, &errors.Object{
			Id:     "c847970d-f142-46da-a70d-37cb2441ad74",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unsupported version.",
			Meta: map[string]any{
				"version": in.Version,
			},
		}
	case "", "v2":
		// v2 is the default version.
		_, err := dom.ProcessDataSourceDirV2(r, processDataSourceDirInput)
		if err != nil {
			return nil, errors.Forward(err, "e7f79fb3-d511-4859-848b-15533364e6f6")
		}
	default:
		return nil, &errors.Object{
			Id:     "8d775b8a-177f-4ce0-8f7e-b03b54a14642",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Unknown version.",
			Meta: map[string]any{
				"version": in.Version,
			},
		}
	}

	out := &ProcessDataSourceOutput{}

	return out, nil
}

type LockInput struct {
	PartnerId  uuid.UUID
	LockerName string
}

type LockOutput struct {
	ExtendCancel context.CancelFunc
	LockerWg     *sync.WaitGroup
}

func (dom *domain) Lock(r *arc.Request, in *LockInput) (*LockOutput, error) {
	switch {
	case dom.config.File.DistributedLockers == nil:
		return nil, &errors.Object{
			Id:     "209c0342-321f-424f-88ab-43428382f385",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing config.File.DistributedLockers.",
		}
	case len(dom.config.File.DistributedLockers) == 0:
		return nil, &errors.Object{
			Id:     "dcb6e952-ec1f-4add-97cd-b28cac9f319f",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Empty config.File.DistributedLockers.",
		}
	}

	distLocker, ok := dom.config.File.DistributedLockers[in.LockerName]
	if !ok {
		return nil, &errors.Object{
			Id:     "7d2f4c23-5572-4348-9f8e-b57ba7f6e851",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing locker.",
		}
	}

	if distLocker.TableName == "" {
		return nil, &errors.Object{
			Id:     "c5757be4-06a0-4fe9-a1b6-bed2b862f6ee",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing locker table name.",
		}
	}

	if distLocker.Keys == nil {
		return nil, &errors.Object{
			Id:     "9e0787ac-9cd1-4ab5-9ae5-aff721b701c6",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing locker keys.",
		}
	}

	partnerName, err := PartnerNameById(in.PartnerId)
	if err != nil {
		return nil, errors.Forward(err, "04e6ab83-941a-4196-aec9-bb75598fb7f7")
	}

	partnerKey, ok := distLocker.Keys[partnerName]
	if !ok {
		return nil, &errors.Object{
			Id:     "658b3897-fdfc-4c5c-a243-1f93cfc0e1be",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing locker partner key.",
		}
	}

	locker := &distsync.DynamoDB{
		PollInterval: time.Second,
		Client:       dynamodb.NewFromConfig(dom.config.AWS.Get("default")),
		TableName:    distLocker.TableName,
	}

	ttl := time.Minute

	lock := &distsync.Lock{
		Id:        partnerKey,
		Type:      distsync.LockTypeWrite,
		Ttl:       ttl,
		NoPolling: true,
		NoReaders: true,
	}

	lockCtx, lockCancel := context.WithTimeout(r.Context(), time.Minute)
	defer lockCancel()

	if err := locker.Lock(lockCtx, lock); err != nil {
		return nil, errors.Forward(err, "0328467c-d538-4b4d-b324-76e064335e90")
	}

	log.Info().
		Str("lock_key", partnerKey).
		Msg("Lock acquired.")

	extendCtx, extendCancel := context.WithCancel(r.Context())

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		// Extend the lock while the worker is running.

		ticker := time.NewTicker(ttl / 2)

		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-extendCtx.Done():
				log.Info().Msg("Extend lock done. Unlocking.")

				if err := locker.Unlock(context.Background()); err != nil {
					log.Error().
						Err(errors.Forward(err, "2fad4b77-d060-4be2-9cd1-83b213bdf690")).
						Send()
				}

				return
			case <-ticker.C:
				if err := locker.Extend(extendCtx); err != nil {
					// Info log because canceling extend is expected.

					log.Info().
						Err(errors.Forward(err, "8c60c639-b460-47e3-b1de-e969a8de6120")).
						Send()

					return
				}
			}
		}
	}()

	out := &LockOutput{
		ExtendCancel: extendCancel,
		LockerWg:     wg,
	}

	return out, nil
}
