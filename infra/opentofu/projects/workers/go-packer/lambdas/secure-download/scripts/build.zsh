source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

LAMBDA_GO_OUT=${ABODEMINE_BUILD_TMP}/bin/bootstrap
echo "LAMBDA_GO_OUT=${LAMBDA_GO_OUT}" >> ${ABODEMINE_BUILD_ENV}

################################################################################
# Go.
################################################################################

# Used by modules/lambda/lambda.tf.
# The binary name MUST be `bootstrap` for AWS Lambda.
# Reference: https://github.com/aws/aws-lambda-go?tab=readme-ov-file#building-your-function.
GO_OUT=${LAMBDA_GO_OUT} \
${MAKE} -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/workers/packer/lambdas/secure-download build

################################################################################
# Config.
################################################################################

mkdir -p ${ABODEMINE_BUILD_TMP}/etc

gomplate \
    -d "env=${ABODEMINE_BUILD_ENV}?type=application/x-env" \
    -d "params=file:${ABODEMINE_BUILD_PARAMS}" \
    -f ${ABODEMINE_BUILD_DIR}/etc/config.yaml.gotmpl \
> ${ABODEMINE_BUILD_TMP}/etc/config.yaml
