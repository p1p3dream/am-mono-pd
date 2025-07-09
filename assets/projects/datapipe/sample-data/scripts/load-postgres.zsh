# This script fetches sample data from S3
# and loads it into the datapipe database (PostgreSQL).

mkdir -p .tmp

dotenv -e ${ABODEMINE_WORKSPACE}/.env.${ABODEMINE_NAMESPACE} \
aws s3 sync \
	s3://testing-mono-overlay-c5tfj3sl/assets/projects/datapipe/sample-data/.tmp/ \
	${ABODEMINE_WORKSPACE}/assets/projects/datapipe/sample-data/.tmp/

for table in \
    "ad_df_assessor" \
    "ad_df_listing" \
    "ad_df_recorder" \
    "ad_df_rental_avm" \
    "ad_geom" \
    "addresses" \
    "current_listings" \
    "fa_df_address" \
    "fa_df_assessor" \
    "fa_df_avm_power" \
    "fa_geom" \
    "fips" \
    "mls_providers" \
    "properties" \
; do
    echo "Processing table ${table}."
    psql -U abodemine -h postgres.abodemine.local -d datapipe -c "truncate ${table}"

    zstd -cd -T0 .tmp/${table}.txt.zst \
    | psql -U abodemine -h postgres.abodemine.local -d datapipe \
        -c "copy ${table} from stdin"
done
