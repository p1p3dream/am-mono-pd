ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_PROJECT_SUFFIX := go-packer
ABODEMINE_PROJECT_SLUG := workers-$(ABODEMINE_PROJECT_SUFFIX)

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)

build: env
	@echo "Nothing to build."

package: env
	@echo "Nothing to package."

configure: env
	@echo "Nothing to configure."

release: env
	zsh -eo pipefail scripts/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
