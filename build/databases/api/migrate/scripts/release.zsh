source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

CONFIG_DIR=$(yq '.config_dir.'${ABODEMINE_NAMESPACE} config.yaml)
MIGRATIONS_DIR=$(yq '.migrations_dir.'${ABODEMINE_NAMESPACE} config.yaml)

ssh ${ABODEMINE_BASTION_SSH} -- mkdir -p ${CONFIG_DIR} ${MIGRATIONS_DIR}

rsync -aP ${ABODEMINE_BUILD_TMP}/etc/sql-migrate-config.yaml ${ABODEMINE_BASTION_SSH}:${CONFIG_DIR}/sql-migrate-config.yaml
rsync -aP --delete ${ABODEMINE_WORKSPACE}/code/sql/databases/${ABODEMINE_DATABASE_NAME}/migrations/ ${ABODEMINE_BASTION_SSH}:${MIGRATIONS_DIR}/
