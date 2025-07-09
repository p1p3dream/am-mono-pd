package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"abodemine/domains/address"
	"abodemine/domains/arc"
	"abodemine/entities"
	"abodemine/lib/consts"
	"abodemine/lib/distsync"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/val"
	"abodemine/projects/datapipe/domains/partners/abodemine"
	"abodemine/repositories/opensearch"
)

type SyncPropertiesInput struct {
	NoLock              bool
	OpenSearchIndexName string
	Version             string
	WorkerId            *uuid.UUID
}

type SyncPropertiesOutput struct{}

func (dom *domain) SyncProperties(r *arc.Request, in *SyncPropertiesInput) (*SyncPropertiesOutput, error) {
	if !in.NoLock {
		lockOut, err := dom.Lock(r, &LockInput{
			PartnerId:  abodemine.PartnerId,
			LockerName: "synther",
		})
		if err != nil {
			return nil, errors.Forward(err, "6ccca780-2a80-437a-8de6-5b4f63f8f50f")
		}

		defer func() {
			lockOut.ExtendCancel()
			lockOut.LockerWg.Wait()
		}()
	}

	log.Info().Msg("Syncing properties.")

	// Ensure we have up-to-date zip5s.
	if err := dom.addressDomain.UpdateZip5Table(r); err != nil {
		return nil, errors.Forward(err, "8fdb92d2-9792-4a46-8929-f0ac15876601")
	}

	selectZip5Out, err := dom.addressDomain.SelectZip5(r, &address.SelectZip5Input{
		OrderBy: "zip5",
	})
	if err != nil {
		return nil, errors.Forward(err, "ec6e908c-a93a-4acc-b40f-15bd17f40c70")
	}

	log.Info().
		Int("zip5_count", len(selectZip5Out.Models)).
		Send()

	batchSize := 1000
	g, gctx := errgroup.WithContext(context.Background())
	g.SetLimit(6)

	r = r.Clone(arc.CloneRequestWithContext(gctx))

	for i, zip5 := range selectZip5Out.Models {
		g.Go(func() error {
			log.Info().
				Int("remaining", len(selectZip5Out.Models)-i).
				Str("zip5", zip5.Zip5).
				Msg("Processing Zip5.")

			for j := 0; ; j++ {
				syncOut, err := dom.SyncPropertiesByZip5(r, &SyncPropertiesByZip5Input{
					BatchCount:          j + 1,
					BatchSize:           batchSize,
					Index:               i,
					OpenSearchIndexName: in.OpenSearchIndexName,
					Total:               len(selectZip5Out.Models),
					Zip5:                zip5.Zip5,
				})
				if err != nil {
					return errors.Forward(err, "69154d85-0265-49f7-9774-ff6e989449ce")
				}

				if syncOut.MatchCount < batchSize {
					break
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Forward(err, "c6bd1860-d59a-492f-98f3-a6852b8d6d88")
	}

	out := &SyncPropertiesOutput{}

	return out, nil
}

type SyncPropertiesByZip5Input struct {
	BatchCount          int
	BatchSize           int
	Index               int
	OpenSearchIndexName string
	Total               int
	Zip5                string
}

type SyncPropertiesByZip5Output struct {
	MatchCount int
}

func (dom *domain) SyncPropertiesByZip5(r *arc.Request, in *SyncPropertiesByZip5Input) (*SyncPropertiesByZip5Output, error) {
	batchId, err := val.NewUUID4()
	if err != nil {
		return nil, errors.Forward(err, "0e171b04-a63d-4119-a373-b7d6239ed6bc")
	}

	sql := `
		with ad_asr as (
			select
				max(attomid) as attomid,
				upper(property_address_full) as address
			from ad_df_assessor
			where
				property_address_zip = $1
				and property_address_full is not null
				and attomid <> 999999999
				and not exists (
					select 1
					from properties
					where properties.ad_attom_id = ad_df_assessor.attomid
				)
			group by address
		), fa_asr as (
			select
				max(property_id) as property_id,
				upper(situs_full_street_address) as address
			from fa_df_assessor
			where
				situs_zip5 = $1
				and situs_full_street_address is not null
				and not exists (
					select 1
					from properties
					where properties.fa_property_id = fa_df_assessor.property_id
				)
			group by address
		), matched_asr as (
			select
				ad_asr.attomid,
				fa_asr.property_id
			from ad_asr
			join fa_asr using (address)
			limit $3
		), new_ad_geom as (
			insert into ad_geom (
				id,
				created_at,
				updated_at,
				attom_id,
				location,
				location_3857
			)
			select
				gen_random_uuid(),
				now(),
				now(),
				matched_asr.attomid,
				st_setsrid(
					st_makepoint(
						ad_df_assessor.property_longitude,
						ad_df_assessor.property_latitude
					),
					4326
				) as location,
				st_transform(
					st_setsrid(
						st_makepoint(
							ad_df_assessor.property_longitude,
							ad_df_assessor.property_latitude
						),
						4326
					),
					3857
				) as location_3857
			from matched_asr
			join ad_df_assessor on matched_asr.attomid = ad_df_assessor.attomid
			on conflict (attom_id)
			do update set
				updated_at = excluded.updated_at,
				location = excluded.location,
				location_3857 = excluded.location_3857
		), new_fa_geom as (
			insert into fa_geom (
				id,
				created_at,
				updated_at,
				property_id,
				location,
				location_3857
			)
			select
				gen_random_uuid(),
				now(),
				now(),
				matched_asr.property_id,
				st_setsrid(
					st_makepoint(
						situs_longitude,
						situs_latitude
					),
					4326
				) as location,
				st_transform(
					st_setsrid(
						st_makepoint(
							situs_longitude,
							situs_latitude
						),
						4326
					),
					3857
				) as location_3857
			from matched_asr
			join fa_df_assessor on matched_asr.property_id = fa_df_assessor.property_id
			on conflict (property_id)
			do update set
				updated_at = excluded.updated_at,
				location = excluded.location,
				location_3857 = excluded.location_3857
		), new_properties as (
			insert into properties (
				id,
				created_at,
				updated_at,
				meta,
				ad_attom_id,
				fa_property_id
			)
			select
				gen_random_uuid() as id,
				now() as created_at,
				now() as updated_at,
				jsonb_build_object(
					'batch_id',
					$2::text
				) as meta,
				matched_asr.attomid,
				matched_asr.property_id
			from matched_asr
			returning *
		), new_addresses as (
			insert into addresses (
				id,
				created_at,
				updated_at,
				meta,
				data_source,
				city,
				county,
				fips,
				full_street_address,
				house_number,
				state,
				street_name,
				street_pos_direction,
				street_pre_direction,
				street_suffix,
				unit_nbr,
				unit_type,
				zip5
			)
			select
				gen_random_uuid() as id,
				now() as created_at,
				now() as updated_at,
				jsonb_build_object(
					'batch_id',
					$2::text,
					'property_id',
					new_properties.id
				) as meta,
				'fa_df_assessor' as data_source,
				normalize_address(fa_df_assessor.situs_city),
				(select county from fips where fips = fa_df_assessor.fips) as county,
				fa_df_assessor.fips,
				normalize_address(fa_df_assessor.situs_full_street_address) as full_street_address,
				fa_df_assessor.situs_house_nbr as house_number,
				fa_df_assessor.situs_state,
				normalize_address(fa_df_assessor.situs_street) as street_name,
				fa_df_assessor.situs_direction_right as street_pos_direction,
				fa_df_assessor.situs_direction_left as street_pre_direction,
				initcap(fa_df_assessor.situs_mode) as street_suffix,
				fa_df_assessor.situs_unit_nbr,
				initcap(fa_df_assessor.situs_unit_type),
				fa_df_assessor.situs_zip5
			from new_properties
			join fa_df_assessor on new_properties.fa_property_id = fa_df_assessor.property_id
			returning *
		)
		select
			new_addresses.city,
			new_addresses.county,
			new_addresses.fips,
			new_addresses.full_street_address,
			new_addresses.house_number,
			new_addresses.id,
			new_addresses.meta->>'property_id',
			new_addresses.state,
			new_addresses.street_name,
			new_addresses.street_pos_direction,
			new_addresses.street_pre_direction,
			new_addresses.street_suffix,
			new_addresses.unit_nbr,
			new_addresses.unit_type,
			new_addresses.updated_at,
			new_addresses.zip5
		from new_addresses
	`

	args := []any{
		in.Zip5,
		batchId,
		in.BatchSize,
	}

	pgxPool, err := r.Dom().SelectPgxPool(consts.ConfigKeyPostgresDatapipe)
	if err != nil {
		return nil, errors.Forward(err, "0c1e1d20-86db-492a-9b09-8842b97992c8")
	}

	tx, err := pgxPool.Begin(r.Context())
	if err != nil {
		return nil, &errors.Object{
			Id:     "34bd3a18-8ddc-4fe4-8b89-45d237990199",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to begin transaction.",
			Cause:  err.Error(),
		}
	}

	defer extutils.RollbackPgxTx(r.Context(), tx, "e1ed52d1-bc58-415a-b3fd-f60acb50061e")

	r = r.Clone(arc.CloneRequestWithPgxTx(consts.ConfigKeyPostgresDatapipe, tx))

	rows, err := extutils.PgxQuery(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "81a4e604-c854-4d32-b00d-f2a53cf9a478")
	}

	defer rows.Close()

	newAddresses := []opensearch.Document{}
	newAddressesIds := []uuid.UUID{}

	for rows.Next() {
		record := &entities.PropertyAddress{}

		if err := rows.Scan(
			&record.City,
			&record.County,
			&record.Fips,
			&record.FullStreetAddress,
			&record.HouseNumber,
			&record.Id,
			&record.Aupid,
			&record.State,
			&record.StreetName,
			&record.StreetPostDirection,
			&record.StreetPreDirection,
			&record.StreetSuffix,
			&record.UnitNumber,
			&record.UnitType,
			&record.UpdatedAt,
			&record.Zip5,
		); err != nil {
			return nil, errors.Forward(err, "4ba7a86c-8c84-4363-a338-8b287b7e9b73")
		}

		newAddresses = append(newAddresses, record)
		newAddressesIds = append(newAddressesIds, *record.Id)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, &errors.Object{
			Id:     "571d3afa-e2d3-4447-84e3-17020ba869a4",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to query rows.",
			Cause:  err.Error(),
		}
	}

	out := &SyncPropertiesByZip5Output{}

	log.Info().
		Int("batch_count", in.BatchCount).
		Str("batch_id", batchId.String()).
		Int("match_count", len(newAddressesIds)).
		Str("zip5", in.Zip5).
		Send()

	if len(newAddressesIds) == 0 {
		return out, nil
	}

	// Update properties with the new addresses.
	// This can't be done in the same query as the insert
	// because of transaction isolation levels.
	sql = `
		with created_addresses as (
			select
				id,
				(meta->>'property_id')::uuid as property_id
			from addresses
			where id = any ($1)
		)
		update properties
		set address_id = created_addresses.id
		from created_addresses
		where properties.id = created_addresses.property_id
	`

	args = []any{
		newAddressesIds,
	}

	_, err = extutils.PgxExec(r, consts.ConfigKeyPostgresDatapipe, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "7856d0fb-6af7-46af-bfda-ee8ab7f03cb0")
	}

	backoff := &distsync.Backoff{
		InitialInterval: time.Second,
		MaxInterval:     5 * time.Minute,
		MaxRetries:      20,
	}

RETRY:
	if _, err := dom.osSearchRepository.PutDocument(r, &opensearch.PutDocumentInput{
		IndexName: in.OpenSearchIndexName,
		Items:     newAddresses,
	}); err != nil {
		firstErr := errors.First(err)

		if firstErr.Code == errors.Code_RESOURCE_EXHAUSTED {
			duration, err := backoff.Next()
			if err != nil {
				return nil, errors.Forward(err, "d63ac939-f215-4d2d-855b-071559662e80")
			}

			log.Warn().
				Int("batch_count", in.BatchCount).
				Str("batch_id", batchId.String()).
				Dur("duration", duration).
				Int("retries", backoff.Retries()).
				Str("zip5", in.Zip5).
				Msg("Sleeping before retry due to resource exhaustion.")

			time.Sleep(duration)

			goto RETRY
		}

		return nil, errors.Forward(err, "e3412fca-5c5b-49c7-9624-c82e967a5f4b")
	}

	if err := tx.Commit(r.Context()); err != nil {
		return nil, &errors.Object{
			Id:     "f964fe6d-4d58-406e-93fc-a77b40a47fcd",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to commit transaction.",
			Cause:  err.Error(),
		}
	}

	out.MatchCount = len(newAddressesIds)

	return out, nil
}
