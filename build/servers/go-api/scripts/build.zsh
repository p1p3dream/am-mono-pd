source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

PROJECT_NAME_UPPERCASE=$(printf "${ABODEMINE_PROJECT_NAME}" | tr '[:lower:]' '[:upper:]')

################################################################################
# Go.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/bin

GO_OUT=${ABODEMINE_BUILD_TMP}/docker/bin/server \
${MAKE} -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/servers/${ABODEMINE_PROJECT_NAME} build

################################################################################
# Certs.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc/ssl

sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-ca.sh abodemine-ca ${ABODEMINE_BUILD_TMP}/docker/etc/ssl
sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-cert.sh abodemine-ca peer abodemine.internal ${ABODEMINE_BUILD_TMP}/docker/etc/ssl

################################################################################
# Config.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc

# Copy casbin.
mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc/casbin/api
rsync -aP ${ABODEMINE_WORKSPACE}/code/casbin/projects/api/ ${ABODEMINE_BUILD_TMP}/docker/etc/casbin/api/

# Copy valkey.
mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/api
rsync -aP ${ABODEMINE_WORKSPACE}/code/lua/valkey/api/ ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/api/

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/session
rsync -aP ${ABODEMINE_WORKSPACE}/code/lua/valkey/session/ ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/session/

mkdir -p ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/token
rsync -aP ${ABODEMINE_WORKSPACE}/code/lua/valkey/token/ ${ABODEMINE_BUILD_TMP}/docker/etc/valkey/token/

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

################################################################################
# Dockerfile.
################################################################################

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -f docker/Dockerfile.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/Dockerfile
