ABODEMINE_TOOL_NAME := csvmine

GO_OUT ?= ${ABODEMINE_WORKSPACE}/.local/build/tools/bin/$(ABODEMINE_TOOL_NAME)
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
		abodemine/tools/$(ABODEMINE_TOOL_NAME)

run:
	dotenv -e ${ABODEMINE_WORKSPACE}/.env.${ABODEMINE_NAMESPACE} \
	go run abodemine/tools/$(ABODEMINE_TOOL_NAME) $(RUN_ARGS)
