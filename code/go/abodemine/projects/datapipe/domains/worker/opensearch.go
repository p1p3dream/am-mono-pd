package worker

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeebo/xxh3"
	"golang.org/x/sync/errgroup"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/lib/distsync"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/domains/partners/abodemine"
	"abodemine/projects/datapipe/entities"
	"abodemine/repositories/opensearch"
)

type LoadOpenSearchInput struct {
	BatchSize      int
	FileBufferSize int
	IndexName      string
	NoLock         bool
	PartnerId      uuid.UUID
	Version        string
}

type LoadOpenSearchOutput struct{}

func (dom *domain) LoadOpenSearch(r *arc.Request, in *LoadOpenSearchInput) (*LoadOpenSearchOutput, error) {
	if !in.NoLock {
		lockOut, err := dom.Lock(r, &LockInput{
			PartnerId:  abodemine.PartnerId,
			LockerName: "osloader",
		})
		if err != nil {
			return nil, errors.Forward(err, "7e8bf527-2dad-4be5-a200-d9dc689b7b25")
		}

		defer func() {
			lockOut.ExtendCancel()
			lockOut.LockerWg.Wait()
		}()
	}

	indexName := strings.TrimSpace(in.IndexName)

	if indexName == "" {
		return nil, &errors.Object{
			Id:     "9549854d-9e25-47f9-b133-61ca6186f015",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "IndexName is required.",
		}
	}

	indexExistsOut, err := dom.osSearchRepository.IndexExists(r, &opensearch.IndexExistsInput{
		IndexName: indexName,
	})
	if err != nil {
		return nil, errors.Forward(err, "e780e0e1-4a95-4170-96ef-fe749dd52202")
	}

	if !indexExistsOut.Exists {
		log.Info().
			Str("index_name", indexName).
			Msg("Creating index.")

		_, err = dom.osSearchRepository.CreateIndex(r, &opensearch.CreateIndexInput{
			IndexName: indexName,
			Body: &opensearch.IndexCreateBody{
				Settings: map[string]*opensearch.IndexCreateBodySetting{
					"index": {
						NumberOfShards:   1,
						NumberOfReplicas: 1,
					},
				},
				Mappings: map[string]any{
					"dynamic": "strict",
					"properties": map[string]*opensearch.IndexCreateBodyMappingProperty{
						"aupid": {
							Type: "keyword",
						},
						"city": {
							Type: "text",
						},
						"county": {
							Type: "text",
						},
						"fips": {
							Type: "keyword",
						},
						"fullStreetAddress": {
							Type: "text",
						},
						"houseNumber": {
							Type: "text",
						},
						"id": {
							Type: "keyword",
						},
						"state": {
							Type: "text",
						},
						"stateFullName": {
							Type: "text",
						},
						"streetName": {
							Type: "text",
						},
						"streetPostDirection": {
							Type: "text",
						},
						"streetPreDirection": {
							Type: "text",
						},
						"streetSuffix": {
							Type: "text",
						},
						"unitNumber": {
							Type: "text",
						},
						"unitType": {
							Type: "text",
						},
						"updatedAt": {
							Type: "date",
						},
						"zip5": {
							Type: "keyword",
						},
					},
				},
			},
		})
		if err != nil {
			return nil, errors.Forward(err, "c0b75aaa-12f7-45f7-9b00-f199d01d76c9")
		}
	}

	indexDirectoryPath := path.Join(
		"opensearch:",
		dom.osSearchRepository.ConfigKey(),
		in.IndexName,
	)

	selectDirectoryOut, err := dom.repository.SelectDataFileDirectoryRecord(r, &SelectDataFileDirectoryRecordInput{
		PartnerId: in.PartnerId,
		Path:      indexDirectoryPath,
	})
	if err != nil {
		return nil, errors.Forward(err, "5e402543-d0ab-44fa-8638-6d2f7a48e1d1")
	}

	directory := selectDirectoryOut.Record

	if directory == nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "5f941c49-8dd7-4614-9453-6db278cea644")
		}

		now := time.Now()

		insertDirectoryOut, err := dom.repository.InsertDataFileDirectoryRecord(r, &InsertDataFileDirectoryRecordInput{
			Record: &entities.DataFileDirectory{
				Id:        id,
				CreatedAt: now,
				UpdatedAt: now,
				PartnerId: in.PartnerId,
				Status:    entities.DataFileDirectoryStatusToDo,
				Path:      indexDirectoryPath,
				Name:      in.IndexName,
			},
		})
		if err != nil {
			return nil, errors.Forward(err, "77d82d64-9c66-42fb-b5bc-70a687c0b692")
		}

		directory = insertDirectoryOut.Record
	}

	out := &LoadOpenSearchOutput{}

	if directory.Status == entities.DataFileDirectoryStatusDone ||
		directory.Status == entities.DataFileDirectoryStatusIgnored {
		log.Debug().
			Str("path", directory.Path).
			Msg("Directory is already processed.")
		return out, nil
	}

	updateDirectoryOut, err := dom.repository.UpdateDataFileDirectoryRecord(r, &UpdateDataFileDirectoryRecordInput{
		Id:        directory.Id,
		UpdatedAt: time.Now(),
		Status:    entities.DataFileDirectoryStatusInProgress,
	})
	if err != nil {
		return nil, errors.Forward(err, "e73892c1-fea0-4f95-8193-c6fbe004c878")
	}

	directory = updateDirectoryOut.Record

	selectFipsOut, err := dom.addressDomain.SelectFips(r, &address.SelectFipsInput{
		OrderBy: "fips",
	})
	if err != nil {
		return nil, errors.Forward(err, "b99bd608-194a-42e5-ac23-8be9acbf5d76")
	}

	log.Info().
		Str("directory_path", directory.Path).
		Int("fips_count", len(selectFipsOut.Models)).
		Send()

	backoff := &distsync.Backoff{
		InitialInterval: time.Second,
		MaxInterval:     5 * time.Minute,
		MaxRetries:      20,
	}

	processedFips := val.NewCache[string, uuid.UUID]()

RETRY:
	g, gctx := errgroup.WithContext(context.Background())
	g.SetLimit(in.FileBufferSize)

	r = r.Clone(arc.CloneRequestWithContext(gctx))

	for _, fips := range selectFipsOut.Models {
		lastId, ok := processedFips.Select(fips.Fips)

		if ok && lastId == uuid.Nil {
			continue
		}

		g.Go(func() error {
			loadObjectOut, err := dom.LoadOpenSearchObject(r, &LoadOpenSearchObjectInput{
				BatchSize:    in.BatchSize,
				DataFileType: abodemine.DataFileTypeSearchAddress,
				DirectoryId:  &directory.Id,
				Fips:         fips.Fips,
				IndexName:    in.IndexName,
				LastId:       lastId,
			})
			if err != nil {
				if loadObjectOut != nil && loadObjectOut.Retry {
					processedFips.Set(fips.Fips, loadObjectOut.LastId)
				}

				return errors.Forward(err, "304be433-3c39-425b-8ce4-29d2bef281ef")
			}

			// Reset processed fips to prevent reprocessing.
			processedFips.Set(fips.Fips, uuid.Nil)

			// Reset backoff after successful operations.
			backoff.Reset()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		firstError := errors.First(err)

		if firstError.Code == errors.Code_RESOURCE_EXHAUSTED {
			duration, err := backoff.Next()
			if err != nil {
				return nil, errors.Forward(err, "115b74e8-916f-4e9c-823d-2c3d6b2cdf53")
			}

			log.Warn().
				Dur("duration", duration).
				Int("retries", backoff.Retries()).
				Msg("Sleeping before retry due to resource exhaustion.")

			time.Sleep(duration)

			goto RETRY
		}

		return nil, errors.Forward(err, "67a8bdbe-5207-4c72-a32a-9182ff2d60c0")
	}

	return out, nil
}

type LoadOpenSearchObjectInput struct {
	BatchSize    int
	DataFileType entities.DataFileType
	DirectoryId  *uuid.UUID
	Fips         string
	IndexName    string
	LastId       uuid.UUID
}

type LoadOpenSearchObjectOutput struct {
	LastId uuid.UUID
	Retry  bool
}

func (dom *domain) LoadOpenSearchObject(r *arc.Request, in *LoadOpenSearchObjectInput) (*LoadOpenSearchObjectOutput, error) {
	if r.Context().Err() != nil {
		return nil, &errors.Object{
			Id:     "8ec5e308-589c-4175-a3be-227856aa502e",
			Code:   errors.Code_CANCELED,
			Detail: "Request cancelled.",
		}
	}

	if in.Fips == "" {
		return nil, &errors.Object{
			Id:     "9e289d64-63b8-4f18-a201-00dda33ff9ea",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Fips is required.",
		}
	}

	log.Info().
		Str("fips", in.Fips).
		Msg("Loading OpenSearch object.")

	hash := val.ByteArray16ToSlice(xxh3.HashString128(in.Fips).Bytes())

	selectObjectOut, err := dom.repository.SelectDataFileObjectRecord(r, &SelectDataFileObjectRecordInput{
		DirectoryId: in.DirectoryId,
		FileType:    in.DataFileType,
		Hash:        hash,
	})
	if err != nil {
		return nil, errors.Forward(err, "675bde8a-c4ff-4cd8-9b30-3c9b8b0b5aed")
	}

	dfObject := selectObjectOut.Record

	if dfObject == nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "230398f7-dcff-4f46-bb08-8d865b22f130")
		}

		now := time.Now()

		insertObjectOut, err := dom.repository.InsertDataFileObjectRecord(r, &InsertDataFileObjectRecordInput{
			Record: &entities.DataFileObject{
				Id:          id,
				CreatedAt:   now,
				UpdatedAt:   now,
				DirectoryId: in.DirectoryId,
				FileType:    in.DataFileType,
				Hash:        hash,
				Status:      entities.DataFileObjectStatusToDo,
				FileName:    in.Fips,
			},
		})
		if err != nil {
			return nil, errors.Forward(err, "4654ccae-fdc9-4182-af46-8735f145b547")
		}

		dfObject = insertObjectOut.Record
	}

	out := &LoadOpenSearchObjectOutput{}

	if dfObject.Status == entities.DataFileObjectStatusDone ||
		dfObject.Status == entities.DataFileObjectStatusIgnored {
		log.Info().
			Str("path", in.Fips).
			Msg("Object already processed.")
		return out, nil
	}

	updateObjectOut, err := dom.repository.UpdateDataFileObjectRecord(r, &UpdateDataFileObjectRecordInput{
		Id:        dfObject.Id,
		UpdatedAt: time.Now(),
		Status:    entities.DataFileObjectStatusInProgress,
	})
	if err != nil {
		return nil, errors.Forward(err, "0623e692-3645-4d7b-bc54-f82afa46d6b9")
	}

	dfObject = updateObjectOut.Record
	lastId := in.LastId

	// Ensure the record count restarts at zero.
	// We will refresh previous records, if any, but don't have to worry
	// about partial batch writes. This is a safe operation.
	dfObject.RecordCount = 0

	for {
		selectAddressesOut, err := dom.addressDomain.SelectPropertyAddress(r, &address.SelectPropertyAddressInput{
			Fips: in.Fips,
			IdGt: &lastId,
			Columns: []string{
				"addresses.city",
				"addresses.county",
				"addresses.fips",
				"addresses.full_street_address",
				"addresses.house_number",
				"addresses.id",
				"addresses.meta->>'property_id'",
				"addresses.state",
				"addresses.street_name",
				"addresses.street_pos_direction",
				"addresses.street_pre_direction",
				"addresses.street_suffix",
				"addresses.unit_nbr",
				"addresses.unit_type",
				"addresses.updated_at",
				"addresses.zip5",
			},
			Limit:                    uint64(in.BatchSize),
			OrderBy:                  []string{"addresses.id"},
			IncludeStateFullName:     true,
			ReturnOpenSearchDocument: true,
		})
		if err != nil {
			return nil, errors.Forward(err, "72308f6c-82f8-47da-95cc-6ffeddd29e6a")
		}

		if len(selectAddressesOut.AddressDocuments) == 0 {
			break
		}

		if _, err := dom.osSearchRepository.PutDocument(r, &opensearch.PutDocumentInput{
			IndexName: in.IndexName,
			Items:     selectAddressesOut.AddressDocuments,
		}); err != nil {
			var out *LoadOpenSearchObjectOutput

			firstErr := errors.First(err)

			if firstErr.Code == errors.Code_RESOURCE_EXHAUSTED {
				out = &LoadOpenSearchObjectOutput{
					LastId: lastId,
					Retry:  true,
				}
			}

			return out, errors.Forward(err, "612166ea-a0b6-4c8e-b024-424afbedc4ea")
		}

		updateObjectOut, err := dom.repository.UpdateDataFileObjectRecord(r, &UpdateDataFileObjectRecordInput{
			Id:          dfObject.Id,
			UpdatedAt:   time.Now(),
			RecordCount: dfObject.RecordCount + int32(len(selectAddressesOut.AddressDocuments)),
		})
		if err != nil {
			return nil, errors.Forward(err, "c44aba29-362e-4194-a524-91357e431d94")
		}

		dfObject = updateObjectOut.Record

		lastId, err = val.UUIDFromString(selectAddressesOut.AddressDocuments[len(selectAddressesOut.AddressDocuments)-1].OpenSearchId())
		if err != nil {
			return nil, errors.Forward(err, "9f9020fb-300d-479a-9664-45dd42c54a8d")
		}
	}

	if _, err := dom.repository.UpdateDataFileObjectRecord(r, &UpdateDataFileObjectRecordInput{
		Id:        dfObject.Id,
		UpdatedAt: time.Now(),
		Status:    entities.DataFileObjectStatusDone,
	}); err != nil {
		return nil, errors.Forward(err, "b78e1f8a-90e8-4fc1-8a12-668bfb4aaad0")
	}

	return out, nil
}
