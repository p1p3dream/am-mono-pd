include $(ABODEMINE_WORKSPACE)/code/make/build.mk

# https://github.com/hairyhenderson/gomplate/releases.
GOMPLATE_VERSION := 4.3.0

# Also update infra/docker/projects/am-env/Makefile.
PACKER_VERSION ?= 0.4.3

ifeq ($(UNAME_M),aarch64)
GOMPLATE_DOWNLOAD_URL := https://github.com/hairyhenderson/gomplate/releases/download/v$(GOMPLATE_VERSION)/gomplate_linux-arm64
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-aarch64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-aarch64
else ifeq ($(UNAME_M),arm64)
GOMPLATE_DOWNLOAD_URL := https://github.com/hairyhenderson/gomplate/releases/download/v$(GOMPLATE_VERSION)/gomplate_linux-arm64
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-aarch64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-aarch64
else ifeq ($(UNAME_M),x86_64)
GOMPLATE_DOWNLOAD_URL := https://github.com/hairyhenderson/gomplate/releases/download/v$(GOMPLATE_VERSION)/gomplate_linux-amd64
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-x86_64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-x86_64
endif

all:;

env: $(ABODEMINE_BUILD_ENV)
	@echo "GOMPLATE_DOWNLOAD_URL=$(GOMPLATE_DOWNLOAD_URL)" >> $(ABODEMINE_BUILD_ENV)
	@echo "PACKER_PATH=$(PACKER_PATH)" >> $(ABODEMINE_BUILD_ENV)
	@echo "PACKER_TARGET=$(PACKER_TARGET)" >> $(ABODEMINE_BUILD_ENV)

build: env
	bash -eo pipefail scripts/build.bash -e "$(ABODEMINE_BUILD_ENV)"

package: env
	bash -eo pipefail scripts/package.bash -e "$(ABODEMINE_BUILD_ENV)"
