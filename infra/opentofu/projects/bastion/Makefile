ABODEMINE_BUILD_PARAMS_VERSION := v2

include $(ABODEMINE_WORKSPACE)/code/make/build.mk

ABODEMINE_CONFIG := ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/bastion/config.${ABODEMINE_NAMESPACE}.yaml
ABODEMINE_PROJECT_SLUG := bastion
ABODEMINE_TOFU_DEPLOYMENT ?= ${ABODEMINE_NAMESPACE}

all:;

env: $(ABODEMINE_BUILD_ENV) $(ABODEMINE_BUILD_PARAMS)
	@echo "ABODEMINE_CONFIG=$(ABODEMINE_CONFIG)" >> $(ABODEMINE_BUILD_ENV)

files/%: $(ABODEMINE_CONFIG) env files/%.gotmpl
	gomplate \
		-d config=$(ABODEMINE_CONFIG) \
		-d "env=$(ABODEMINE_BUILD_ENV)?type=application/x-env" \
		-d params=$(ABODEMINE_BUILD_PARAMS) \
		-f $@.gotmpl \
	> $@

FILES_ETC_DEPS := \
	files/etc/systemd/system/abodemine-init.service \
	files/etc/systemd/system/os-alpha-pipe.service \
	files/etc/systemd/system/pg-alpha-pipe.service \
	files/etc/systemd/system/pg-beta-pipe.service \
	files/etc/systemd/system/vk-alpha-pipe.service

files/etc.tar.zst: env $(FILES_ETC_DEPS)
	zsh -eo pipefail scripts/create-etc-tar.zsh -e "$(ABODEMINE_BUILD_ENV)"

FILES_HOME_DEPS := \
	files/home/.gitconfig \
	files/home/.zshrc \
	files/home/.cargo/config.toml \
	files/home/.abodemine/bin/abodemine-init.zsh \
	files/home/.abodemine/bin/add-portfwd-user.zsh \
	files/home/.abodemine/bin/run-actions-runner.zsh \
	files/home/.abodemine/code/zsh/base.zsh \
	files/home/.ssh/config \
	files/home/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519 \
	files/home/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519.pub

files/home/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519 files/home/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519.pub:
	ssh-keygen -t ed25519 -f files/home/.ssh/abodeminebot-${ABODEMINE_NAMESPACE}.ed25519 -N "" -C "abodeminebot"

files/home.tar.zst: env $(FILES_HOME_DEPS)
	zsh -eo pipefail scripts/create-home-tar.zsh -e "$(ABODEMINE_BUILD_ENV)"

user_data.bash: $(ABODEMINE_CONFIG) env files/etc.tar.zst files/home.tar.zst user_data.bash.gotmpl
	gomplate \
		-d config=$(ABODEMINE_CONFIG) \
		-d "env=$(ABODEMINE_BUILD_ENV)?type=application/x-env" \
		-f $@.gotmpl \
	> $@

user_data.bash.gz: user_data.bash
	gzip -fk $^

build: env user_data.bash.gz

package: env
	@echo "Nothing to package."

configure: env $(OPENTOFU_DEPS)
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/configure.zsh -e "$(ABODEMINE_BUILD_ENV)"

release: env
	zsh -eo pipefail $(ABODEMINE_WORKSPACE)/code/sh/infra/opentofu/release.zsh -e "$(ABODEMINE_BUILD_ENV)"
