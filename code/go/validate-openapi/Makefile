.PHONY: build validate

build:
	go mod tidy
	go mod download
	go build -o validate-openapi validate.go
	chmod +x validate-openapi

validate: build
	./validate-openapi ../../../../../docs/projects/api/abodemine-openapi-3.0.0.yaml

clean:
	rm -f validate-openapi 