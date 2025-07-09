source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

CONFIG_OUT_DIR=${ABODEMINE_BUILD_TMP}/etc
CONFIG_OUT_FILE=${CONFIG_OUT_DIR}/config.yaml.zst
SECURE_DOWNLOAD_BUCKET=$(yq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".s3_buckets."secure-download".name' ${ABODEMINE_BUILD_PARAMS})

mkdir -p ${CONFIG_OUT_DIR}

zstd -T0 -6 \
    -o ${CONFIG_OUT_FILE} \
    ${ABODEMINE_BUILD_DIR}/etc/config.yaml

aws s3 cp \
    ${CONFIG_OUT_FILE} \
    s3://${SECURE_DOWNLOAD_BUCKET}/packer/config.yaml.zst
