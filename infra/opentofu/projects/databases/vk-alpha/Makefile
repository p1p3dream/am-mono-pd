ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_SUFFIX := vk-alpha
ABODEMINE_PROJECT_SLUG := databases-$(ABODEMINE_PROJECT_SUFFIX)
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)

build: env
	@echo "Nothing to build."

package: env
	@echo "Nothing to package."

configure: env $(OPENTOFU_DEPS)
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
