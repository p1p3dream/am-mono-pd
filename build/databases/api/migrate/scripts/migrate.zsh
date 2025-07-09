source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

CONFIG_DIR=$(yq '.config_dir.'${ABODEMINE_NAMESPACE} config.yaml)

ssh ${ABODEMINE_BASTION_SSH} -- \
    /opt/abodemine/pkg/sql-migrate/sql-migrate $1 \
        -config ${CONFIG_DIR}/sql-migrate-config.yaml \
        -env=${ABODEMINE_NAMESPACE}
