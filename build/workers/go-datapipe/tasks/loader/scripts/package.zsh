source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

TASK_KEY=$(printf "projects/%s/tasks/%s" "${ABODEMINE_PROJECT_SLUG}" "${ABODEMINE_TASK_NAME}")

IMAGE_REPOSITORY=$(jq -r '."'${TASK_KEY}'".ecr_repository' ${ABODEMINE_BUILD_PARAMS})
IMAGE_TAG="${IMAGE_REPOSITORY}:${ABODEMINE_BUILD_ID}"

docker buildx build \
    --tag ${IMAGE_TAG} \
    ${ABODEMINE_BUILD_TMP}/docker

echo "ABODEMINE_APP_IMAGE_TAG=${IMAGE_TAG}" >> ${ABODEMINE_BUILD_ENV}
