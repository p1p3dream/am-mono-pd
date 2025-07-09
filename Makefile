ABODEMINE_BUILD_PARAMS ?= .config/params.${ABODEMINE_NAMESPACE}.json

all:;

aws-sso-login:
	dotenv -e .env.${ABODEMINE_NAMESPACE} \
	aws sso login --no-browser --use-device-code

params:
	mkdir -p .config
	dotenv -e .env.${ABODEMINE_NAMESPACE} \
	code/sh/bin/aws-ssm-to-json.zsh $(ABODEMINE_BUILD_PARAMS_VERSION) \
	> $(ABODEMINE_BUILD_PARAMS)
