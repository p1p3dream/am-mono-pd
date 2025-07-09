ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_NAME := datapipe
ABODEMINE_PROJECT_SUFFIX := go-$(ABODEMINE_PROJECT_NAME)
ABODEMINE_PROJECT_SLUG := workers-$(ABODEMINE_PROJECT_SUFFIX)
ABODEMINE_TASK_NAME := loader
ABODEMINE_TASK_SLUG := $(ABODEMINE_PROJECT_SLUG)-$(ABODEMINE_TASK_NAME)
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)
	@echo "GOARCH=$(GOARCH)" >> $(ABODEMINE_BUILD_ENV)

build: env
	zsh -eo pipefail scripts/build.zsh -e "$(ABODEMINE_BUILD_ENV)"

package: env
	zsh -eo pipefail scripts/package.zsh -e "$(ABODEMINE_BUILD_ENV)"

configure: env $(OPENTOFU_DEPS)
	OPENTOFU_DIR=opentofu \
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail scripts/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
