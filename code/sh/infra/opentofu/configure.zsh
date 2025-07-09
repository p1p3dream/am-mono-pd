source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

OPENTOFU_SRC=${ABODEMINE_BUILD_DIR}
OPENTOFU_DST=${ABODEMINE_BUILD_TMP}/opentofu

if [[ -n "${OPENTOFU_DIR}" ]]; then
    OPENTOFU_SRC+=/${OPENTOFU_DIR}
fi

find ${OPENTOFU_SRC} -type d | xargs -I{} tofu -chdir={} fmt

mkdir -p ${OPENTOFU_DST}

exclude=()

if [[ -n "${RSYNC_EXCLUDE}" ]]; then
    IFS=',' read -rA exclude_items <<< "${RSYNC_EXCLUDE}"
    for item in "${exclude_items[@]}"; do
        exclude+=("--exclude=${item}")
    done
fi

rsync -a "${exclude[@]}" ${OPENTOFU_SRC}/ ${OPENTOFU_DST}/

args=()

if [[ -n "${ABODEMINE_CONFIG}" ]]; then
    args+=(-d "config=${ABODEMINE_CONFIG}")
fi

args+=(
    -d "params=${ABODEMINE_BUILD_PARAMS}"
    -f "${OPENTOFU_DST}/config.auto.tfvars.gotmpl"
)

gomplate "${args[@]}" > "${OPENTOFU_DST}/config.auto.tfvars"
