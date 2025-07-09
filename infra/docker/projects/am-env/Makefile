all:;

include $(ABODEMINE_WORKSPACE)/code/make/core.mk

ABODEMINE_SHELL_SSH_SERVER_ENABLED ?= true

DOCKER ?= docker
DOCKER_COMPOSE ?= docker-compose
DOCKER_PROJECT ?= $(shell basename $(CURDIR))
DOCKER_COMPOSE_OPTIONS ?= --project-name $(DOCKER_PROJECT) --file compose.yaml

# https://hub.docker.com/r/hairyhenderson/gomplate/tags.
GOMPLATE_IMAGE ?= hairyhenderson/gomplate:v4.3.0
ifeq ($(ABODEMINE_SHELL),true)
GOMPLATE := gomplate
else
GOMPLATE := $(DOCKER) run --rm --volume $(CURDIR):/works --workdir /works $(GOMPLATE_IMAGE)
endif

# Also update build/tools/github-runner/Makefile.
PACKER_VERSION ?= 0.4.3

ifeq ($(UNAME_M),aarch64)
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-aarch64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-aarch64
else ifeq ($(UNAME_M),arm64)
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-aarch64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-aarch64
else ifeq ($(UNAME_M),x86_64)
PACKER_PATH := packer/v$(PACKER_VERSION)/packer-$(PACKER_VERSION)-x86_64-unknown-linux-gnu.tar.zst
PACKER_TARGET := linux-x86_64
endif

# https://hub.docker.com/r/mikefarah/yq/tags.
YQ_IMAGE ?= mikefarah/yq:4.44.6
ifeq ($(ABODEMINE_SHELL),true)
YQ := yq
else
YQ := $(DOCKER) run --rm --volume $(CURDIR):/works --workdir /works $(YQ_IMAGE)
endif

# config.yaml MUST be first because they are used by other templates.
GENERATED_BUILD_FILE_DEPS := \
	config.yaml \
	compose.yaml \
	services/shell/.config/rclone/rclone.conf \
	services/shell/.gnupg/gpg-agent.conf \
	services/shell/.ssh/config \
	services/shell/home/.gitconfig \
	services/shell/home/.tmux.conf \
	services/shell/home/.zshrc \
	services/shell/Dockerfile \
	services/shell/entrypoint.sh \
	services/tools/Dockerfile

BUILD_FILE_DEPS := \
	$(GENERATED_BUILD_FILE_DEPS)

# Create an empty local configuration file if it does not exist.
.env:
	if [ ! -f $@ ]; then \
		touch $@; \
	fi

compose.yaml: config.yaml compose.yaml.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

# Create an empty local configuration file if it does not exist.
config.local.yaml:
	if [ ! -f $@ ]; then \
		touch $@; \
	fi

# Merge default and local configurations, where local is overlaid on top of default.
config.yaml: config.default.yaml config.local.yaml
	$(YQ) \
		-M \
		-P \
		eval-all \
		'. as $$item ireduce ({}; . * $$item )' \
		$^ \
	> $@

services/%/Dockerfile: .env config.yaml services/%/Dockerfile.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $@.gotmpl \
	> $@

services/shell/.cargo/config.toml: config.yaml services/shell/.cargo/config.toml.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/.config/rclone/rclone.conf: config.yaml services/shell/.config/rclone/rclone.conf.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/.gnupg/gpg-agent.conf: config.yaml services/shell/.gnupg/gpg-agent.conf.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/.ssh/config: config.yaml services/shell/.ssh/config.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

ifeq ($(ABODEMINE_SHELL_SSH_SERVER_ENABLED),true)
BUILD_FILE_DEPS += services/shell/.ssh/id_ed25519 services/shell/.ssh/id_ed25519.pub

services/shell/.ssh/id_ed25519 services/shell/.ssh/id_ed25519.pub:
	ssh-keygen -t ed25519 -f $@ -N ""
endif

services/shell/home/.gitconfig: config.yaml services/shell/home/.gitconfig.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/home/.tmux.conf: config.yaml services/shell/home/.tmux.conf.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/home/.zshrc: config.yaml services/shell/home/.zshrc.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/shell/entrypoint.sh: config.yaml services/shell/entrypoint.sh.gotmpl
	$(GOMPLATE) \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

build: $(BUILD_FILE_DEPS)
	mkdir -p \
		services/shell/.cache \
		services/shell/.cache/bin \
		services/shell/.cache/etc \
		services/tools/.cache/bin \
		services/tools/.cache/downloads \
		services/tools/.cache/etc

	$(DOCKER) build \
		--tag $(DOCKER_PROJECT)-tools \
		services/tools

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	curl \
		-L \
		-o services/tools/.cache/downloads/packer.tar.zst \
		-H "Authorization: Bearer $$(grep '^ABODEMINE_SECURE_DOWNLOAD_TOKEN=' .env | sed 's/ABODEMINE_SECURE_DOWNLOAD_TOKEN=//')" \
		-H "X-AbodeMine-S3-Object: $(PACKER_PATH)" \
		$$(grep '^ABODEMINE_SECURE_DOWNLOAD_ENDPOINT=' .env | sed 's/ABODEMINE_SECURE_DOWNLOAD_ENDPOINT=//')

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	zstd -cd -T0 services/tools/.cache/downloads/packer.tar.zst \
	| $(DOCKER) run --interactive --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	bsdtar -C services/tools/.cache/bin -xf -

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	curl \
		-L \
		-o services/tools/.cache/downloads/config.yaml.zst \
		-H "Authorization: Bearer $$(grep '^ABODEMINE_SECURE_DOWNLOAD_TOKEN=' .env | sed 's/ABODEMINE_SECURE_DOWNLOAD_TOKEN=//')" \
		-H "X-AbodeMine-S3-Object: packer/config.yaml.zst" \
		$$(grep '^ABODEMINE_SECURE_DOWNLOAD_ENDPOINT=' .env | sed 's/ABODEMINE_SECURE_DOWNLOAD_ENDPOINT=//')

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	zstd -o services/tools/.cache/etc/config.yaml -df services/tools/.cache/downloads/config.yaml.zst

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	services/tools/.cache/bin/packer --version

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	services/tools/.cache/bin/packer \
		--tracing-level info \
		--config services/tools/.cache/etc/config.yaml \
		--downloads-dir services/shell/.cache/downloads \
		--target $(PACKER_TARGET) \
		--packages-dir /opt/abodemine/pkg \
		--temp-dir /tmp/abodemine/packer \
		--seven-zip 7za \
		--profiles default,am-env \
		--redownload-transient \
		--download-only \

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	rsync -a services/tools/.cache/bin/packer services/shell/.cache/bin/packer

	$(DOCKER) run --rm --volume $(CURDIR):/works $(DOCKER_PROJECT)-tools \
	rsync -a services/tools/.cache/etc/config.yaml services/shell/.cache/etc/config.yaml

	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) build

network: config.yaml
	DOCKER_NETWORK_NAME=$$($(YQ) -M '.vars.docker.network.name' config.yaml); \
	if ! $(DOCKER) network ls --filter name=$${DOCKER_NETWORK_NAME} --format {{.Name}} | grep -q "^$${DOCKER_NETWORK_NAME}$$"; then \
		$(DOCKER) network create $${DOCKER_NETWORK_NAME}; \
	fi

setup: compose.yaml network build

clean: down
	DOCKER_NETWORK_NAME=$$($(YQ) -M '.vars.docker.network.name' config.yaml); \
	if $(DOCKER) network ls --filter name=$${DOCKER_NETWORK_NAME} --format {{.Name}} | grep -q "^$${DOCKER_NETWORK_NAME}$$"; then \
		$(DOCKER) network rm $${DOCKER_NETWORK_NAME}; \
	fi
	for item in \
		$(DOCKER_PROJECT)_works \
	; do \
		if $(DOCKER) volume ls --filter name=$${item} --format {{.Name}} | grep -q "^$${item}$$"; then \
			$(DOCKER) volume rm $${item}; \
		fi; \
	done

	rm -f $(GENERATED_BUILD_FILE_DEPS)
	rm -rf services/shell/.cache
	rm -rf services/tools/.cache

down: compose.yaml
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) down --remove-orphans

restart:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) restart

start:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) start

stop:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) stop

up: setup down
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) up --detach

logs:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) logs --follow

stats:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) stats

sync:
	ABODEMINE_USER=$$($(YQ) -M '.user.email' config.yaml | sed 's/@.*//'); \
	rsync -a \
		--exclude=\.git \
		--exclude=services/shell/.cache/downloads \
		--delete \
		$${ABODEMINE_USER}@am-env-shell:/works/src/mono/infra/docker/projects/am-env/ \
		${ABODEMINE_WORKSPACE}/infra/docker/projects/am-env/; \
	rsync -a \
		$${ABODEMINE_USER}@am-env-shell:/works/src/mono/code/make/core.mk \
		${ABODEMINE_WORKSPACE}/code/make/core.mk

mount-nfs: config.yaml
	sudo mkdir -p /Volumes/abodemine/am-works
	sudo mount \
		-t nfs -o noacl,noatime,port=$$($(YQ) -M '.services.shell.vars.nfs.port' config.yaml),vers=4 \
		127.0.0.1:/ \
		/Volumes/abodemine/am-works

	sudo mkdir -p /Volumes/abodemine/am-works-src
	sudo mount \
		-t nfs -o noacl,noatime,port=$$($(YQ) -M '.services.shell.vars.nfs.port' config.yaml),vers=4 \
		127.0.0.1:/src \
		/Volumes/abodemine/am-works-src

#	Work in progress to use automount.
#	echo "AmWorksSrc -fstype=nfs,vers=4,port=$$($(YQ) -M '.services.shell.vars.nfs.port' config.yaml) 127.0.0.1:/src" \
#	| sudo tee /etc/auto_abodemine
#	sudo sed -I.bkp '/auto_abodemine/ d' /etc/auto_master
#	echo "/Volumes/abodemine auto_abodemine" | sudo tee -a /etc/auto_master
#	sudo automount -vc

umount-nfs:
	sudo umount /Volumes/abodemine/am-works
	sudo umount /Volumes/abodemine/am-works-src
