DOC_TEMPLATE = docs/snippets/config.snippet
.PHONY: run documentation

build: 
	@echo "Building binary..."
	go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go

test:
	@echo "Running Golang tests..."
	go test ./...

doc:
	@echo "Generating documentation..."
	go doc -u configuration.Config | sed -n '/type Config struct {/,/}/p' > ${DOC_TEMPLATE}
	@cat ${DOC_TEMPLATE}

docker-build:
	@echo "Building docker image with tag 'dev'..."
	docker build -t jetspotter:dev .

docker-run: docker-build
	@echo "Running docker container with tag 'dev'..."
	docker run jetspotter:dev
