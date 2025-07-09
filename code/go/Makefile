DEPS := \
	go.work

all: $(DEPS)

check-unique-uuids: check-unique-uuids-v1

check-unique-uuids-v1:
	@UUIDS=$$(\
		grep \
			--include=\*.go \
			--exclude=\*_test.go \
			-EI \
			-nor \
			"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}" \
			. \
			| sed 's/:/\t/g' \
			| sort -k 3 \
			| uniq -D -f 2\
			) && \
		if [ ! -z "$$UUIDS" ]; then \
			echo "Error: duplicate UUIDs found:"; \
			printf "$$UUIDS\n" && \
			exit 1; \
		fi

ci: check-unique-uuids
	$(MAKE) -C abodemine ci

go.work:
	go work init

sync-workspace: go.work
# Ensure each package is tidy.
	for modfile in $(shell find . -type f -name go.mod); do \
		dir=$$(dirname $$modfile); \
		cd $$dir && go mod tidy && cd $(CURDIR); \
	done

# Add each package to the workspace.
	for modfile in $(shell find . -type f -name go.mod); do \
		dir=$$(dirname $$modfile); \
		go work use $$dir; \
	done

	go work sync
