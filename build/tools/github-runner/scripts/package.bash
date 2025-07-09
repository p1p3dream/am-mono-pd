source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

docker buildx build \
    --tag github-runner:latest \
    ${ABODEMINE_BUILD_TMP}/docker
