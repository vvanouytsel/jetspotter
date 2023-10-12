CONFIG_SNIPPET = docs/snippets/config.snippet
OUTPUT_SNIPPET = docs/snippets/output.snippet
OVERVIEW_SNIPPET = docs/snippets/overview.snippet
.PHONY: run documentation

build: 
	@echo "Building binary..."
	go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go

test:
	@echo "Running Golang tests..."
	go test ./...

doc:
	@echo "Generating documentation..."
	go doc -u configuration.Config | sed -n '/type Config struct {/,/}/p' > ${CONFIG_SNIPPET}
	@cat ${CONFIG_SNIPPET}
	go doc -u jetspotter.AircraftOutput | sed -n '/type AircraftOutput struct {/,/}/p' > ${OUTPUT_SNIPPET}
	@cat ${OUTPUT_SNIPPET}
	@echo "# Overview\n" > helm/jetspotter/README.md
	@cat ${OVERVIEW_SNIPPET} >> helm/jetspotter/README.md
	@cat helm/jetspotter/README.md

docker-build:
	@echo "Building docker image with tag 'dev'..."
	docker build -t ghcr.io/vvanouytsel/jetspotter:dev .

docker-run: docker-build
	@echo "Running docker container with tag 'dev'..."
	docker run ghcr.io/vvanouytsel/jetspotter:dev
