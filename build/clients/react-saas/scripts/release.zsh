source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

DISTRIBUTION_BUCKET=$(yq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".distribution_bucket' ${ABODEMINE_BUILD_PARAMS})

# Everything except index.html and meta files because we don't want
# to update them before their requirements are ready.
aws s3 sync \
    --exclude "index.html" \
    --exclude "robots.txt" \
    --exclude "manifest.json" \
    ${ABODEMINE_BUILD_TMP}/www/ \
    s3://${DISTRIBUTION_BUCKET}/

# Everything ignored before.
aws s3 sync \
    --exclude "*" \
    --include "index.html" \
    --include "robots.txt" \
    --include "manifest.json" \
    ${ABODEMINE_BUILD_TMP}/www/ \
    s3://${DISTRIBUTION_BUCKET}/
