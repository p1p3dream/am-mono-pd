all:;

DOCKER ?= docker
DOCKER_COMPOSE ?= docker compose
DOCKER_PROJECT ?= $(shell basename $(CURDIR))
DOCKER_PROJECT_PROFILE_ALL ?= all
DOCKER_PROJECT_PROFILE_DEFAULT ?= default
DOCKER_PROJECT_RELPATH := $(shell realpath --relative-to=$(ABODEMINE_WORKSPACE) $(CURDIR))
DOCKER_COMPOSE_OPTIONS ?= --project-name $(DOCKER_PROJECT) --file compose.yaml

POSTGRES_SUPERUSER := postgres
POSTGRES_DB_API ?= api
POSTGRES_DB_DATAPIPE ?= datapipe
POSTGRES_DB_SAAS ?= saas
POSTGRES_PSQL_ARGS ?= -U $(POSTGRES_SUPERUSER)
POSTGRES_USER ?= abodemine

# config.yaml MUST be first because they are used by other templates.
BUILD_FILE_DEPS := \
	config.yaml \
	compose.yaml \
	services/dc-go/Dockerfile \
	services/dc-python/Dockerfile \
	services/dc-root/Dockerfile \
	services/dc-typescript/Dockerfile \
	services/envoy/etc/clusters.yaml \
	services/envoy/etc/filter-chains.yaml \
	services/envoy/etc/main.yaml \
	services/nginx/etc/servers/admin.conf \
	services/nginx/etc/servers/saas-whitelabel.conf \
	services/nginx/etc/servers/saas.conf \
	services/nginx/etc/nginx.conf \
	services/opensearch/certs/root-ca.pem \
	services/opensearch/certs/chain.pem \
	services/opensearch/certs/key.pem \
	services/opensearch/Dockerfile \
	services/postgres/Dockerfile

# Create an empty local configuration file if it does not exist.
.env:
	if [ ! -f $@ ]; then \
		touch $@; \
	fi

.PHONY: certs
certs:
	$(MAKE) -C $(ABODEMINE_WORKSPACE)/etc/ssl

compose.yaml: .env config.yaml compose.yaml.gotmpl
	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $@.gotmpl \
	| yq -MP \
	> $@

# Create an empty local configuration file if it does not exist.
config.local.yaml:
	if [ ! -f $@ ]; then \
		touch $@; \
	fi

# Merge default and local configurations, where local is overlaid on top of default.
config.yaml: config.default.yaml config.local.yaml
	yq \
		-M \
		-P \
		eval-all \
		'. as $$item ireduce ({}; . * $$item )' \
		$^ \
	| yq '.vars.docker.project.name = "$(DOCKER_PROJECT)"' \
	| yq '.vars.docker.project.relpath = "$(DOCKER_PROJECT_RELPATH)"' \
	> $@

services/%/Dockerfile: config.yaml services/%/Dockerfile.gotmpl
	gomplate \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/envoy/etc/clusters.yaml: config.yaml services/envoy/etc/clusters.yaml.gotmpl
	gomplate \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/envoy/etc/filter-chains.yaml: config.yaml services/envoy/etc/filter-chains.yaml.gotmpl
	gomplate \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/envoy/etc/main.yaml: config.yaml services/envoy/etc/clusters.yaml services/envoy/etc/filter-chains.yaml services/envoy/etc/main.yaml.gotmpl
	gomplate \
		-d config=config.yaml \
		--template clusters=services/envoy/etc/clusters.yaml \
		--template filter-chains=services/envoy/etc/filter-chains.yaml \
		-f $@.gotmpl \
	> $@

services/nginx/etc/servers/%.conf: config.yaml services/nginx/etc/servers/%.conf.gotmpl
	gomplate \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/nginx/etc/nginx.conf: config.yaml services/nginx/etc/nginx.conf.gotmpl
	gomplate \
		-d config=config.yaml \
		-f $@.gotmpl \
	> $@

services/opensearch/certs/root-ca.pem: certs
	mkdir -p $(dir $@)
	cp $(ABODEMINE_WORKSPACE)/etc/ssl/abodemine-ca.pem $@

services/opensearch/certs/chain.pem: certs
	mkdir -p $(dir $@)
	cp $(ABODEMINE_WORKSPACE)/etc/ssl/abodemine.local-peer-chain.pem $@

services/opensearch/certs/key.pem: certs
	mkdir -p $(dir $@)
	cp $(ABODEMINE_WORKSPACE)/etc/ssl/abodemine.local-peer-key.pem $@

build: $(BUILD_FILE_DEPS)
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $(DOCKER_PROJECT_PROFILE_ALL) build

restart-postgres:
ifneq ($(shell $(DOCKER) ps --format json | jq -r '. | select(.Names == "$(DOCKER_PROJECT)-postgres-1")'),)
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) restart postgres
else
	-$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) stop postgres
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) run --detach --rm --name postgres postgres
endif
	@printf "Waiting for postgres to restart..."
	@while ! nmap -Pn postgres -p 5432 | grep -q open; do \
		printf "."; \
		sleep 1; \
	done; \
	sleep 1; \
	echo "OK"

stop-standalone-postgres:
ifeq ($(shell $(DOCKER) ps --format json | jq -r '. | select(.Names == "$(DOCKER_PROJECT)-postgres-1")'),)
	$(DOCKER) stop postgres
endif

db/drop: restart-postgres
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_API)
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_DATAPIPE)
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_SAAS)

db/create:
	if [ "$$($(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres psql $(POSTGRES_PSQL_ARGS) -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$(POSTGRES_USER)';")" != "1" ]; then \
		$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createuser $(POSTGRES_PSQL_ARGS) -DRs $(POSTGRES_USER); \
	fi

	if [ "$$($(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres psql $(POSTGRES_PSQL_ARGS) -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DB_API)';")" != "1" ]; then \
		$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_API); \
	fi

	if [ "$$($(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres psql $(POSTGRES_PSQL_ARGS) -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DB_DATAPIPE)';")" != "1" ]; then \
		$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_DATAPIPE); \
	fi
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_DATAPIPE) -c "create extension if not exists postgis"
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_DATAPIPE) -c "create extension if not exists postgis_topology"

	if [ "$$($(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres psql $(POSTGRES_PSQL_ARGS) -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$(POSTGRES_DB_SAAS)';")" != "1" ]; then \
		$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_SAAS); \
	fi

test-db/drop:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_API)_test
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_DATAPIPE)_test
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres dropdb $(POSTGRES_PSQL_ARGS) --if-exists $(POSTGRES_DB_SAAS)_test

test-db/create: test-db/drop db/create migrate/up
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_API)_test
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_DATAPIPE)_test
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres createdb $(POSTGRES_PSQL_ARGS) -O $(POSTGRES_USER) $(POSTGRES_DB_SAAS)_test
# Copy schema only.
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -s -d $(POSTGRES_DB_API) \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_API)_test -q
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -s -d $(POSTGRES_DB_DATAPIPE) \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_DATAPIPE)_test -q
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -s -d $(POSTGRES_DB_SAAS) \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_SAAS)_test -q
# Copy migrations data.
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_API) -a -t zz_migrations \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_API)_test -q
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_DATAPIPE) -a -t zz_migrations \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_DATAPIPE)_test -q
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec postgres pg_dump $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_SAAS) -a -t zz_migrations \
	| $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) exec -T postgres psql $(POSTGRES_PSQL_ARGS) -d $(POSTGRES_DB_SAAS)_test -q

test-db: restart-postgres test-db/create stop-standalone-postgres

db/recreate: db/drop db

db: test-db

${ABODEMINE_WORKSPACE}/code/sql/databases/api/etc/sql-migrate-config.yaml: \
	config.yaml \
	${ABODEMINE_WORKSPACE}/code/sql/databases/api/etc/sql-migrate-config.yaml.gotmpl

	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $@.gotmpl \
	> $@

migrate-api/%: ${ABODEMINE_WORKSPACE}/code/sql/databases/api/etc/sql-migrate-config.yaml
	sql-migrate $* \
		-config ${ABODEMINE_WORKSPACE}/code/sql/databases/api/etc/sql-migrate-config.yaml \
		-env=development \
		$(DATAPIPE_MIGRATION_NAME)

${ABODEMINE_WORKSPACE}/code/sql/databases/datapipe/etc/sql-migrate-config.yaml: \
	config.yaml \
	${ABODEMINE_WORKSPACE}/code/sql/databases/datapipe/etc/sql-migrate-config.yaml.gotmpl

	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $@.gotmpl \
	> $@

migrate-datapipe/%: ${ABODEMINE_WORKSPACE}/code/sql/databases/datapipe/etc/sql-migrate-config.yaml
	sql-migrate $* \
		-config ${ABODEMINE_WORKSPACE}/code/sql/databases/datapipe/etc/sql-migrate-config.yaml \
		-env=development \
		$(DATAPIPE_MIGRATION_NAME)

${ABODEMINE_WORKSPACE}/code/sql/databases/saas/etc/sql-migrate-config.yaml: \
	config.yaml \
	${ABODEMINE_WORKSPACE}/code/sql/databases/saas/etc/sql-migrate-config.yaml.gotmpl

	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $@.gotmpl \
	> $@

migrate-saas/%: ${ABODEMINE_WORKSPACE}/code/sql/databases/saas/etc/sql-migrate-config.yaml
	sql-migrate $* \
		-config ${ABODEMINE_WORKSPACE}/code/sql/databases/saas/etc/sql-migrate-config.yaml \
		-env=development \
		$(DATAPIPE_MIGRATION_NAME)

migrate/up: migrate-api/up migrate-datapipe/up migrate-saas/up

setup: certs build
	$(MAKE) db

clean: down/all
	for item in \
		$(DOCKER_PROJECT)_postgres \
		$(DOCKER_PROJECT)_valkey \
	; do \
		if $(DOCKER) volume ls --filter name=$${item} --format {{.Name}} | grep -q "^$${item}$$"; then \
			$(DOCKER) volume rm $${item}; \
		fi; \
	done

	rm -f $(BUILD_FILE_DEPS)

down: compose.yaml
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $(DOCKER_PROJECT_PROFILE_DEFAULT) down --remove-orphans

down/%: compose.yaml
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $* down --remove-orphans

restart:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $(DOCKER_PROJECT_PROFILE_DEFAULT) restart $(SERVICE)

restart/%:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $* restart $(SERVICE)

start:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) start

stop:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $(DOCKER_PROJECT_PROFILE_DEFAULT) stop

stop/%:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $* stop

up: setup down conf
	if [ ! -d $(ABODEMINE_WORKSPACE)/code/python/abodemine/.venv ]; then \
		$(MAKE) -C $(ABODEMINE_WORKSPACE)/code/python/abodemine venv-create; \
	fi
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $(DOCKER_PROJECT_PROFILE_ALL) up --detach

up/%: setup
	if [ ! -d $(ABODEMINE_WORKSPACE)/code/python/abodemine/.venv ]; then \
		$(MAKE) -C $(ABODEMINE_WORKSPACE)/code/python/abodemine venv-create; \
	fi
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) --profile $* up --detach

logs:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) logs --follow

stats:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_OPTIONS) stats

CONF_DEPS := \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/api/conf/local.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/api/conf/local-ci.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/ci/conf/local-ci.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/datapipe/conf/local.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/datapipe/conf/local-ci.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/search/conf/local.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/search/conf/local-ci.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/saas/conf/local.yaml \
	$(ABODEMINE_WORKSPACE)/code/go/abodemine/projects/saas/conf/local-ci.yaml

$(ABODEMINE_WORKSPACE)/code/%/local.yaml: .env $(ABODEMINE_WORKSPACE)/code/%/config.yaml.gotmpl
	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $(ABODEMINE_WORKSPACE)/code/$*/config.yaml.gotmpl \
	> $@

$(ABODEMINE_WORKSPACE)/code/%/local-ci.yaml: .env $(ABODEMINE_WORKSPACE)/code/%/config.yaml.gotmpl
	LOCAL_CI=true \
	gomplate \
		-d config=config.yaml \
		-d env=.env?type=application/x-env \
		-f $(ABODEMINE_WORKSPACE)/code/$*/config.yaml.gotmpl \
	> $@

conf: .env config.yaml $(CONF_DEPS)
