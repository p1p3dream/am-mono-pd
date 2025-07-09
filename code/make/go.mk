.PHONY: $(DEPS)

$(DEPS):
	$(MAKE) -C $@

CMD_DIR ?= $(CURDIR)
GO_PKG = $(shell realpath $(CMD_DIR) --relative-to $(ABODEMINE_WORKSPACE)/code/go)

GOOS ?= linux
GOARCH ?= amd64

GO ?= go
CGO_ENABLED := 0

GO_OUT ?= $(CMD_DIR)/go_out

$(GO_OUT):
	$(GO) build \
		-o $(GO_OUT) \
		-ldflags "-s -w" \
		$(GO_PKG)
