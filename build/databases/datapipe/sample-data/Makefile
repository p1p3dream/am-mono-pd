ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_DATABASE_NAME := datapipe
ABODEMINE_DATABASE_SERVER := pg-beta

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)
	@echo "ABODEMINE_BASTION_SSH=$(ABODEMINE_BASTION_SSH)" >> $(ABODEMINE_BUILD_ENV)
	@echo "ABODEMINE_DATABASE_NAME=$(ABODEMINE_DATABASE_NAME)" >> $(ABODEMINE_BUILD_ENV)
	@echo "ABODEMINE_DATABASE_SERVER=$(ABODEMINE_DATABASE_SERVER)" >> $(ABODEMINE_BUILD_ENV)

build: env
	@echo "Nothing to build."

package: env
	@echo "Nothing to package."

configure: env
	zsh -eo pipefail scripts/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail scripts/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
