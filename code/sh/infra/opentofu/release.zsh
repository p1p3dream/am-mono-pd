source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

tofu \
    -chdir=${ABODEMINE_BUILD_TMP}/opentofu \
    init

tofu \
    -chdir=${ABODEMINE_BUILD_TMP}/opentofu \
    apply
