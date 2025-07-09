source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

# Compile the env file.

gomplate \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -d "params=file:${ABODEMINE_BUILD_PARAMS}" \
    -f .env.gotmpl \
> ${ABODEMINE_WORKSPACE}/code/typescript/apps/${ABODEMINE_PROJECT_NAME}/.env.${ABODEMINE_NAMESPACE}

# Ensure pnpm is present.
corepack enable
echo "pnpm: $(pnpm --version)"

# Install deps.
make -C ${ABODEMINE_WORKSPACE}/code/typescript install

# Sync overlay.

OVERLAY_BUCKET=$(yq -r '.aws.buckets."mono-overlay"' ${ABODEMINE_BUILD_PARAMS})
PARAMS_AWS_REGION=$(yq -r '.aws.region' ${ABODEMINE_BUILD_PARAMS})

rclone config update --non-interactive s3 region=${PARAMS_AWS_REGION}
rclone copy \
    s3:${OVERLAY_BUCKET}/code/typescript/apps/${ABODEMINE_PROJECT_NAME}/ \
    ${ABODEMINE_WORKSPACE}/code/typescript/apps/${ABODEMINE_PROJECT_NAME}/ \
    --ignore-existing \
    --progress

# Build the app.

mkdir -p ${ABODEMINE_BUILD_TMP}/www

OUT_DIR=${ABODEMINE_BUILD_TMP}/www \
make -C ${ABODEMINE_WORKSPACE}/code/typescript/apps/${ABODEMINE_PROJECT_NAME} build
