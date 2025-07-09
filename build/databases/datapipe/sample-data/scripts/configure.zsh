source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

BIN_DIR=$(yq '.bin_dir.'${ABODEMINE_NAMESPACE} config.yaml)

ssh ${ABODEMINE_BASTION_SSH} -- \
    mkdir -p ${BIN_DIR}

rsync -aP \
    ${ABODEMINE_BUILD_DIR}/bin/generate-sample-data.zsh \
    ${ABODEMINE_BASTION_SSH}:${BIN_DIR}/generate-sample-data.zsh

rsync -aP \
    --delete \
    ~/.aws/ \
    ${ABODEMINE_BASTION_SSH}:~/.aws/
