ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_NAME := api
ABODEMINE_PROJECT_SUFFIX := go-$(ABODEMINE_PROJECT_NAME)
ABODEMINE_PROJECT_SLUG := servers-$(ABODEMINE_PROJECT_SUFFIX)
ABODEMINE_APP_IMAGE_REPOSITORY := $(ABODEMINE_PROJECT_SLUG)
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)
	@echo "ABODEMINE_API_FLAGS=${ABODEMINE_API_FLAGS}" >> $(ABODEMINE_BUILD_ENV)
	@echo "ABODEMINE_APP_IMAGE_REPOSITORY=$(ABODEMINE_APP_IMAGE_REPOSITORY)" >> $(ABODEMINE_BUILD_ENV)
	@echo "GOARCH=$(GOARCH)" >> $(ABODEMINE_BUILD_ENV)

build: env
	zsh -eo pipefail scripts/build.zsh -e "$(ABODEMINE_BUILD_ENV)"

package: env
	zsh -eo pipefail scripts/package.zsh -e "$(ABODEMINE_BUILD_ENV)"

configure: env $(OPENTOFU_DEPS)
	zsh -eo pipefail scripts/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail scripts/release.zsh -e "$(ABODEMINE_BUILD_ENV)"

wait-deployment: env
	zsh -eo pipefail scripts/wait-deployment.zsh -e "$(ABODEMINE_BUILD_ENV)"

tests: env
	zsh -eo pipefail scripts/tests.zsh -e "$(ABODEMINE_BUILD_ENV)"
