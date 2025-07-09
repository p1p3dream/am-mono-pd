source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

LAMBDA_KEY=$(printf "projects/%s/lambdas/%s" "${ABODEMINE_PROJECT_SLUG}" "${ABODEMINE_LAMBDA_NAME}")
SQS_QUEUE_URL=$(jq -r '."'${LAMBDA_KEY}'".config.sqs_queue_url' ${ABODEMINE_BUILD_PARAMS})

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "attom-data", "task": "fetcher"}' \
#     --message-group-id "attom-data" \
#     --message-deduplication-id "$(uuidgen)"

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "attom-data", "task": "loader"}' \
#     --message-group-id "attom-data" \
#     --message-deduplication-id "$(uuidgen)"

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "first-american", "task": "fetcher"}' \
#     --message-group-id "first-american" \
#     --message-deduplication-id "$(uuidgen)"

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "first-american", "task": "loader"}' \
#     --message-group-id "first-american" \
#     --message-deduplication-id "$(uuidgen)"

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "abodemine", "task": "osloader"}' \
#     --message-group-id "abodemine" \
#     --message-deduplication-id "$(uuidgen)"

# aws sqs send-message \
#     --queue-url "${SQS_QUEUE_URL}" \
#     --message-body '{"partner": "abodemine", "task": "synther"}' \
#     --message-group-id "abodemine" \
#     --message-deduplication-id "$(uuidgen)"
