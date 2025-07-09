source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

DOTENV=${ABODEMINE_WORKSPACE}/infra/docker/projects/am-env/.env

if [[ "${GITHUB_ACTIONS}" == "true" ]]; then
    echo "Downloading gomplate..."
    curl -sL ${GOMPLATE_DOWNLOAD_URL} -o /usr/local/bin/gomplate
    chmod +x /usr/local/bin/gomplate
elif [[ -f "${DOTENV}" ]]; then
    ABODEMINE_SECURE_DOWNLOAD_ENDPOINT=$(rg -No '(ABODEMINE_SECURE_DOWNLOAD_ENDPOINT)=(.*)' -r '$2' ${DOTENV})
    ABODEMINE_SECURE_DOWNLOAD_TOKEN=$(rg -No '(ABODEMINE_SECURE_DOWNLOAD_TOKEN)=(.*)' -r '$2' ${DOTENV})
fi

mkdir -p \
    ${ABODEMINE_BUILD_TMP}/docker/.cache/bin \
    ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads \
    ${ABODEMINE_BUILD_TMP}/docker/.cache/etc \
    ${ABODEMINE_BUILD_TMP}/docker/.config/rclone \
    ${ABODEMINE_BUILD_TMP}/docker/.docker \
    ${ABODEMINE_BUILD_TMP}/docker/home

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -d "env=${ENV_FILE}?type=application/x-env" \
    -f docker/Dockerfile.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/Dockerfile

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -f docker/.config/rclone/rclone.conf.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/.config/rclone/rclone.conf

rsync -a docker/.docker/ ${ABODEMINE_BUILD_TMP}/docker/.docker/

gomplate \
    -d "config=file:${ABODEMINE_BUILD_DIR}/config.yaml" \
    -f docker/home/.zshrc.gotmpl \
> ${ABODEMINE_BUILD_TMP}/docker/home/.zshrc

curl \
    -L \
    -o ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads/packer.tar.zst \
    -H "Authorization: Bearer ${ABODEMINE_SECURE_DOWNLOAD_TOKEN}" \
    -H "X-AbodeMine-S3-Object: ${PACKER_PATH}" \
    ${ABODEMINE_SECURE_DOWNLOAD_ENDPOINT}

zstd -cd -T0 ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads/packer.tar.zst \
| tar -C ${ABODEMINE_BUILD_TMP}/docker/.cache/bin -xf -

curl \
    -L \
    -o ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads/config.yaml.zst \
    -H "Authorization: Bearer ${ABODEMINE_SECURE_DOWNLOAD_TOKEN}" \
    -H "X-AbodeMine-S3-Object: packer/config.yaml.zst" \
    ${ABODEMINE_SECURE_DOWNLOAD_ENDPOINT}

zstd -o ${ABODEMINE_BUILD_TMP}/docker/.cache/etc/config.yaml -df ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads/config.yaml.zst

${ABODEMINE_BUILD_TMP}/docker/.cache/bin/packer --version

${ABODEMINE_BUILD_TMP}/docker/.cache/bin/packer \
    --tracing-level info \
    --config ${ABODEMINE_BUILD_TMP}/docker/.cache/etc/config.yaml \
    --downloads-dir ${ABODEMINE_BUILD_TMP}/docker/.cache/downloads \
    --target ${PACKER_TARGET} \
    --packages-dir /opt/abodemine/pkg \
    --temp-dir /tmp/abodemine/packer \
    --seven-zip 7za \
    --download-only

if [[ "${GITHUB_ACTIONS}" == "true" ]]; then
	echo "ABODEMINE_BUILD_CONTEXT=${ABODEMINE_BUILD_TMP}/docker" >> "${GITHUB_OUTPUT}"
fi
