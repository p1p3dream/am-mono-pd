source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

mkdir -p ${ABODEMINE_BUILD_TMP}/etc

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -d "params=file:${ABODEMINE_BUILD_PARAMS}" \
    -f etc/sql-migrate-config.yaml.gotmpl \
> ${ABODEMINE_BUILD_TMP}/etc/sql-migrate-config.yaml
