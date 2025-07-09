# This script loads datapipe.addresses into the local (dev) OpenSearch instance.

psql -U abodemine -h postgres.abodemine.local -d datapipe -c "delete from data_file_objects where file_type = 299117701"

OPENSEARCH_INDEX=$(yq -r '.values."servers-go-api".string."os-addresses-index"' ${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono/config.yaml)
OPENSEARCH_PASSWORD=$(cat ${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono/.env | grep OPENSEARCH_INITIAL_ADMIN_PASSWORD | cut -d '=' -f2)

curl \
    -X DELETE \
    -u "admin:${OPENSEARCH_PASSWORD}" \
    "https://opensearch.abodemine.local:9200/${OPENSEARCH_INDEX}" \
|| true

make -C \
    ${ABODEMINE_WORKSPACE}/code/go/abodemine/workers/datapipe/tasks/osloader \
    run \
    RUN_ARGS="run --partner-id 7ee2e306-d03f-4f72-abf8-3ac5df4796ab --index-name ${OPENSEARCH_INDEX} --no-lock"
