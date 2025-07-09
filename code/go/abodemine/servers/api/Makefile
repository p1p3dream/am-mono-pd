ABODEMINE_PROJECT_NAME := api

GO_OUT ?= ${ABODEMINE_WORKSPACE}/.local/build/servers-go-$(ABODEMINE_PROJECT_NAME)/bin/server

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
		abodemine/servers/$(ABODEMINE_PROJECT_NAME)
	upx $(GO_OUT)

run:
	go run abodemine/servers/$(ABODEMINE_PROJECT_NAME) --config ${ABODEMINE_WORKSPACE}/code/go/abodemine/projects/$(ABODEMINE_PROJECT_NAME)/conf/local.yaml listen

watch:
	watchexec --restart --workdir ${ABODEMINE_WORKSPACE}/code/go -- \
	$(MAKE) -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/servers/$(ABODEMINE_PROJECT_NAME) run

