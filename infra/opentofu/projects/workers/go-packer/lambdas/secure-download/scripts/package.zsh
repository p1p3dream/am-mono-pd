source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

mkdir -p ${ABODEMINE_BUILD_TMP}/opentofu/modules/lambda

LAMBDA_FUNCTION_PAYLOAD=${ABODEMINE_BUILD_TMP}/opentofu/modules/lambda/lambda_function_payload.zip

7za a ${LAMBDA_FUNCTION_PAYLOAD} \
    ${LAMBDA_GO_OUT} \
    ${ABODEMINE_BUILD_TMP}/etc/config.yaml
