source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

OPENTOFU_SRC=${ABODEMINE_BUILD_DIR}/opentofu
OPENTOFU_DST=${ABODEMINE_BUILD_TMP}/opentofu

find ${OPENTOFU_SRC} -type d | xargs -I{} tofu -chdir={} fmt

mkdir -p ${OPENTOFU_DST}
rsync -a ${OPENTOFU_SRC}/ ${OPENTOFU_DST}/

gomplate \
    -d params=${ABODEMINE_BUILD_PARAMS} \
    -f ${OPENTOFU_DST}/config.auto.tfvars.gotmpl \
> ${OPENTOFU_DST}/config.auto.tfvars
