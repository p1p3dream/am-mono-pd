# This script fetches sample data from S3
# and loads it into the saas database (PostgreSQL).

mkdir -p .tmp

dotenv -e ${ABODEMINE_WORKSPACE}/.env.${ABODEMINE_NAMESPACE} \
aws s3 sync \
	s3://testing-mono-overlay-c5tfj3sl/assets/projects/saas/sample-data/.tmp/ \
	${ABODEMINE_WORKSPACE}/assets/projects/saas/sample-data/.tmp/

psql -U abodemine -h postgres.abodemine.local -d saas -c " \
truncate organizations cascade;
"

for table in \
    "organizations" \
    "roles" \
    "users" \
; do
    echo "Processing table ${table}."

    zstd -cd -T0 .tmp/${table}.txt.zst \
    | psql -U abodemine -h postgres.abodemine.local -d saas \
        -c "copy ${table} from stdin"
done
