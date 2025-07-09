source ${ABODEMINE_WORKSPACE}/code/sh/make.sh

ECS_CLUSTER_ARN=$(jq -r '."aws/ecs".clusters."main-fargate-servers-shared".arn' ${ABODEMINE_BUILD_PARAMS})

if [[ -z ${ECS_CLUSTER_ARN} ]]; then
    echo "ECS_CLUSTER_ARN is required."
    exit 1
fi

ECS_SERVICE_ARN=$(jq -r '."projects/'${ABODEMINE_PROJECT_SLUG}'".ecs.services.main.arn' ${ABODEMINE_BUILD_PARAMS})

if [[ -z ${ECS_SERVICE_ARN} ]]; then
    echo "ECS_SERVICE_ARN is required."
    exit 2
fi

echo "Waiting for ECS service to stabilize..."

aws ecs wait services-stable \
    --cluster ${ECS_CLUSTER_ARN} \
    --services ${ECS_SERVICE_ARN}

# TODO: Create better rollout output with
# aws ecs describe-services --cluster "$ECS_CLUSTER_ARN" --services "$ECS_SERVICE_ARN" --output json.
