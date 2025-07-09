source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

BASE_DIR=${ABODEMINE_BUILD_DIR}/files/etc
OUTPUT_DIR=${ABODEMINE_BUILD_TMP}/$(uuidgen)

mkdir -p \
    ${OUTPUT_DIR}/systemd/system

cp -v ${BASE_DIR}/systemd/system/*.service ${OUTPUT_DIR}/systemd/system/

bsdtar \
    --format ustar \
    --no-fflags \
    -C ${OUTPUT_DIR} \
    -cf - \
    . \
| zstd \
    -T0 \
    -18 \
    -fo ${ABODEMINE_BUILD_DIR}/files/etc.tar.zst

BUCKET=s3://$(yq '.deployment_vars.buckets.mono-build' ${ABODEMINE_CONFIG})

aws s3 cp ${ABODEMINE_BUILD_DIR}/files/etc.tar.zst ${BUCKET}/${ABODEMINE_BUILD_ID}/etc.tar.zst
