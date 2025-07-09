ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_SUFFIX := go-datapipe
ABODEMINE_PROJECT_SLUG := workers-$(ABODEMINE_PROJECT_SUFFIX)
ABODEMINE_LAMBDA_NAME := task-launcher
ABODEMINE_LAMBDA_SLUG := $(ABODEMINE_PROJECT_SLUG)-$(ABODEMINE_LAMBDA_NAME)
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)

build: env
	zsh -eo pipefail scripts/build.zsh -e "$(ABODEMINE_BUILD_ENV)"

package: env
	zsh -eo pipefail scripts/package.zsh -e "$(ABODEMINE_BUILD_ENV)"

configure: env $(OPENTOFU_DEPS)
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/release.zsh -e "$(ABODEMINE_BUILD_ENV)"

# For debugging or testing purposes.
# This is not a standard step in the build process.
send-message: env
	zsh -eo pipefail scripts/send-message.zsh -e "$(ABODEMINE_BUILD_ENV)"
