ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_NAME := saas
ABODEMINE_PROJECT_SUFFIX := react-$(ABODEMINE_PROJECT_NAME)
ABODEMINE_PROJECT_SLUG := clients-$(ABODEMINE_PROJECT_SUFFIX)
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)

build: env
	zsh -eo pipefail scripts/build.zsh -e "$(ABODEMINE_BUILD_ENV)"

package: env
	@echo "Nothing to package."

configure: env $(OPENTOFU_DEPS)
	@echo "Nothing to configure."

release: env
	zsh -eo pipefail scripts/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
