source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

BIN_DIR=$(yq '.bin_dir.'${ABODEMINE_NAMESPACE} config.yaml)

ssh ${ABODEMINE_BASTION_SSH} -- \
    zsh ${BIN_DIR}/generate-sample-data.zsh "datapipe_${ABODEMINE_NAMESPACE}"
