set -eo pipefail

PSQL_SERVICE="$1"
PSQL_COMMAND=(psql service=${PSQL_SERVICE})

rm -rf /works/samples/datapipe
mkdir -p /works/samples/datapipe

echo "Removing previous sample_property_ids table."
${PSQL_COMMAND} -c "drop table if exists sample_property_ids"

echo "Creating sample_property_ids table."
${PSQL_COMMAND} -c "create table sample_property_ids (id uuid)"

echo "Inserting sample property IDs."
${PSQL_COMMAND} -c " \
insert into sample_property_ids
(id)
select id
from properties
where fips = '06037'
limit 100000 \
"

# Ensure we have properties with current_listings.
echo "Inserting sample property IDs with current listings."
${PSQL_COMMAND} -c " \
insert into sample_property_ids
(id)
select properties.id
from properties
join current_listings on properties.ad_attom_id = current_listings.attom_id
where
    properties.fips = '06037'
    and not exists (select 1 from sample_property_ids where id = properties.id)
limit 100000 \
"

echo "Exporting ad_df_assessor data."
${PSQL_COMMAND} -c "\copy (
    select ad_df_assessor.*
    from sample_property_ids
    join properties using (id)
    join ad_df_assessor on properties.ad_attom_id = ad_df_assessor.attomid
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/ad_df_assessor.txt.zst

echo "Exporting ad_df_listing data."
${PSQL_COMMAND} -c "\copy (
    select ad_df_listing.*
    from sample_property_ids
    join properties using (id)
    join ad_df_listing on properties.ad_attom_id = ad_df_listing.attom_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/ad_df_listing.txt.zst

echo "Exporting ad_df_recorder data."
${PSQL_COMMAND} -c "\copy (
    select ad_df_recorder.*
    from sample_property_ids
    join properties using (id)
    join ad_df_recorder on properties.ad_attom_id = ad_df_recorder.attomid
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/ad_df_recorder.txt.zst

echo "Exporting ad_df_rental_avm data."
${PSQL_COMMAND} -c "\copy (
    select ad_df_rental_avm.*
    from sample_property_ids
    join properties using (id)
    join ad_df_rental_avm on properties.ad_attom_id = ad_df_rental_avm.attomid
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/ad_df_rental_avm.txt.zst

echo "Exporting ad_geom data."
${PSQL_COMMAND} -c "\copy (
    select ad_geom.*
    from sample_property_ids
    join properties using (id)
    join ad_geom on properties.ad_attom_id = ad_geom.attom_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/ad_geom.txt.zst

echo "Exporting addresses data."
${PSQL_COMMAND} -c "\copy (
    select addresses.*
    from sample_property_ids
    join properties using (id)
    join addresses on properties.address_id = addresses.id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/addresses.txt.zst

echo "Exporting current_listings data."
${PSQL_COMMAND} -c "\copy (
    select current_listings.*
    from sample_property_ids
    join properties using (id)
    join current_listings on properties.ad_attom_id = current_listings.attom_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/current_listings.txt.zst

echo "Exporting fa_df_address data."
${PSQL_COMMAND} -c "\copy (
    select fa_df_address.*
    from sample_property_ids
    join properties using (id)
    join fa_df_address on properties.fa_address_master_id = fa_df_address.address_master_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/fa_df_address.txt.zst

echo "Exporting fa_df_assessor data."
${PSQL_COMMAND} -c "\copy (
    select fa_df_assessor.*
    from sample_property_ids
    join properties using (id)
    join fa_df_assessor on properties.fa_property_id = fa_df_assessor.property_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/fa_df_assessor.txt.zst

echo "Exporting fa_df_avm_power data."
${PSQL_COMMAND} -c "\copy (
    select fa_df_avm_power.*
    from sample_property_ids
    join properties using (id)
    join fa_df_avm_power on properties.fa_property_id = fa_df_avm_power.property_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/fa_df_avm_power.txt.zst

echo "Exporting fa_geom data."
${PSQL_COMMAND} -c "\copy (
    select fa_geom.*
    from sample_property_ids
    join properties using (id)
    join fa_geom on properties.fa_property_id = fa_geom.property_id
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/fa_geom.txt.zst

echo "Exporting fips data."
${PSQL_COMMAND} -c "\copy (
    select *
    from fips
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/fips.txt.zst

echo "Exporting mls_providers data."
${PSQL_COMMAND} -c "\copy (
    select *
    from mls_providers
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/mls_providers.txt.zst

echo "Exporting properties data."
${PSQL_COMMAND} -c "\copy (
    select properties.*
    from sample_property_ids
    join properties using (id)
) to stdout" \
| zstd -T0 -6 -o /works/samples/datapipe/properties.txt.zst

aws s3 sync \
    /works/samples/datapipe/ \
    s3://testing-mono-overlay-c5tfj3sl/assets/projects/datapipe/sample-data/.tmp/ \
    --exclude "*.txt"
