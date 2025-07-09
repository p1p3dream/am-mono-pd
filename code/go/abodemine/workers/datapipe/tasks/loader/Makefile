ABODEMINE_PROJECT_NAME := datapipe
ABODEMINE_TASK_NAME := loader
ABODEMINE_DATAPIPE_CONFIG_PATH ?= ${ABODEMINE_WORKSPACE}/code/go/abodemine/projects/$(ABODEMINE_PROJECT_NAME)/conf/local.yaml

GO_OUT ?= ${ABODEMINE_WORKSPACE}/.local/build/workers-go-$(ABODEMINE_PROJECT_NAME)/tasks/$(ABODEMINE_TASK_NAME)/bin/worker

# Go env vars.
GOOS ?= linux
GOARCH ?= arm64

build:
	CGO_ENABLED=0 \
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	go build \
		-ldflags " \
			-s \
			-w \
			-X 'abodemine/lib/app.buildId=${ABODEMINE_BUILD_ID}' \
			-X 'abodemine/lib/app.buildVersion=${ABODEMINE_BUILD_VERSION}' \
			" \
		-o $(GO_OUT) \
		abodemine/workers/$(ABODEMINE_PROJECT_NAME)/tasks/$(ABODEMINE_TASK_NAME)
	upx $(GO_OUT)

run:
	go run abodemine/workers/$(ABODEMINE_PROJECT_NAME)/tasks/$(ABODEMINE_TASK_NAME) --config $(ABODEMINE_DATAPIPE_CONFIG_PATH) $(RUN_ARGS)

watch:
	watchexec --restart --workdir ${ABODEMINE_WORKSPACE}/code/go -- \
	$(MAKE) -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/workers/$(ABODEMINE_PROJECT_NAME)/tasks/$(ABODEMINE_TASK_NAME) run $(WATCH_ARGS)
