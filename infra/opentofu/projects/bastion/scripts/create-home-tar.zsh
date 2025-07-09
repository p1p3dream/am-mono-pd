source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

BASE_DIR=${ABODEMINE_BUILD_DIR}/files/home
OUTPUT_DIR=${ABODEMINE_BUILD_TMP}/$(uuidgen)

mkdir -p \
    ${OUTPUT_DIR}/.cargo \
    ${OUTPUT_DIR}/.docker \
    ${OUTPUT_DIR}/.abodemine/bin \
    ${OUTPUT_DIR}/.abodemine/code/zsh \
    ${OUTPUT_DIR}/.ssh

cp -v ${BASE_DIR}/.cargo/config.toml ${OUTPUT_DIR}/.cargo/config.toml
cp -v ${BASE_DIR}/.docker/config.json ${OUTPUT_DIR}/.docker/config.json

cp -v ${BASE_DIR}/.abodemine/bin/*.zsh ${OUTPUT_DIR}/.abodemine/bin/
chmod +x ${OUTPUT_DIR}/.abodemine/bin/*.zsh
cp -v ${BASE_DIR}/.abodemine/code/zsh/*.zsh ${OUTPUT_DIR}/.abodemine/code/zsh/

cp -v ${BASE_DIR}/.ssh/config ${OUTPUT_DIR}/.ssh/config
cat ${BASE_DIR}/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519.pub >> ${OUTPUT_DIR}/.ssh/authorized_keys
chmod 0700 ${OUTPUT_DIR}/.ssh
chmod 0600 ${OUTPUT_DIR}/.ssh/*

cp -v \
    ${BASE_DIR}/.gitconfig \
    ${BASE_DIR}/.inputrc \
    ${BASE_DIR}/.login_conf \
    ${BASE_DIR}/.tmux.conf \
    ${BASE_DIR}/.zshrc \
    ${OUTPUT_DIR}/

bsdtar \
    --format ustar \
    --no-fflags \
    -C ${OUTPUT_DIR} \
    -cf - \
    . \
| zstd \
    -T0 \
    -18 \
    -fo ${ABODEMINE_BUILD_DIR}/files/home.tar.zst

BUCKET=s3://$(yq '.deployment_vars.buckets.mono-build' ${ABODEMINE_CONFIG})

aws s3 cp ${ABODEMINE_BUILD_DIR}/files/home.tar.zst ${BUCKET}/${ABODEMINE_BUILD_ID}/home.tar.zst
