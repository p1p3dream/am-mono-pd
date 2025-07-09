# This script fetches sample data from S3
# and loads it into the api database (PostgreSQL).

mkdir -p .tmp

dotenv -e ${ABODEMINE_WORKSPACE}/.env.${ABODEMINE_NAMESPACE} \
aws s3 sync \
	s3://testing-mono-overlay-c5tfj3sl/assets/projects/api/sample-data/.tmp/ \
	${ABODEMINE_WORKSPACE}/assets/projects/api/sample-data/.tmp/

psql -U abodemine -h postgres.abodemine.local -d api -c " \
truncate clients;
truncate api_keys cascade;
"

for table in \
    "api_keys" \
    "clients" \
; do
    echo "Processing table ${table}."

    zstd -cd -T0 .tmp/${table}.txt.zst \
    | psql -U abodemine -h postgres.abodemine.local -d api \
        -c "copy ${table} from stdin"
done
