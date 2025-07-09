source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

docker push ${ABODEMINE_APP_IMAGE_TAG}

tofu \
    -chdir=${ABODEMINE_BUILD_TMP}/opentofu \
    init

tofu \
    -chdir=${ABODEMINE_BUILD_TMP}/opentofu \
    apply \
    -auto-approve
