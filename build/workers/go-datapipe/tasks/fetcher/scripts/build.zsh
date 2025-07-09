source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

PROJECT_NAME_UPPERCASE=$(printf "${ABODEMINE_PROJECT_NAME}" | tr '[:lower:]' '[:upper:]')

################################################################################
# Go.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/bin

GO_OUT=${ABODEMINE_BUILD_TMP}/docker/bin/worker \
${MAKE} -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/workers/${ABODEMINE_PROJECT_NAME}/tasks/${ABODEMINE_TASK_NAME} build

################################################################################
# Config.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc

# Generate secure key.
ENC_ITER=$((RANDOM % 100000 + 100000))
ENC_PASS="$(openssl rand -base64 32)"

TMP_SUFFIX_1=$(pwgen -A 8 1)

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -d "params=file:${ABODEMINE_BUILD_PARAMS}" \
    -f etc/config.yaml.gotmpl \
> ${ABODEMINE_BUILD_TMP}/.config.yaml.${TMP_SUFFIX_1}

openssl enc -aes-256-cbc -pbkdf2 -iter ${ENC_ITER} -salt \
    -in ${ABODEMINE_BUILD_TMP}/.config.yaml.${TMP_SUFFIX_1} \
    -out ${ABODEMINE_BUILD_TMP}/docker/etc/config.yaml \
    -pass "pass:${ENC_PASS}"

echo "ABODEMINE_${PROJECT_NAME_UPPERCASE}_CONFIG_PATH=/app/etc/config.yaml" >> ${ABODEMINE_BUILD_ENV}
echo "ABODEMINE_${PROJECT_NAME_UPPERCASE}_CONFIG_PATH_ENC_ITER=${ENC_ITER}" >> ${ABODEMINE_BUILD_ENV}
echo "ABODEMINE_${PROJECT_NAME_UPPERCASE}_CONFIG_PATH_ENC_PASS=${ENC_PASS}" >> ${ABODEMINE_BUILD_ENV}

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/home/.config/rclone

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -d "params=file:${ABODEMINE_BUILD_PARAMS}" \
    -f home/.config/rclone/rclone.conf.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/home/.config/rclone/rclone.conf

################################################################################
# Dockerfile.
################################################################################

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -f docker/Dockerfile.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/Dockerfile
