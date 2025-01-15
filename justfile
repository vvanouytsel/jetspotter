config_snippet := "docs/snippets/config.snippet"
output_snippet := "docs/snippets/output.snippet"
overview_snippet := "docs/snippets/overview.snippet"

just:
    just --list

build: 
	go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go

run:
	go run cmd/jetspotter/*
	
test:
	go test ./...

generate-doc:
	echo "Generating documentation..."
	go doc -u configuration.Config | sed -n '/type Config struct {/,/}/p' > {{config_snippet}}
	cat {{config_snippet}}
	go doc -u jetspotter.AircraftOutput | sed -n '/type AircraftOutput struct {/,/}/p' > {{output_snippet}}
	cat {{output_snippet}}
	echo -e "# Overview\n" > helm/jetspotter/README.md
	echo -e "[![GitHub repository](https://img.shields.io/badge/GitHub-jetspotter-green?logo=github)](https://github.com/vvanouytsel/jetspotter)\n" >> helm/jetspotter/README.md
	cat {{overview_snippet}} >> helm/jetspotter/README.md
	cat helm/jetspotter/README.md

build-container:
    #!/bin/bash
    ct=$(command -v podman || command -v docker)
    [[ "$ct" == "" ]] && echo "Please install Podman or Docker." && exit 1
    
    $ct build -t ghcr.io/vvanouytsel/jetspotter:dev .

run-container: build-container
    #!/bin/bash
    ct=$(command -v podman || command -v docker)
    [[ "$ct" == "" ]] && echo "Please install Podman or Docker." && exit 1
    
    $ct run ghcr.io/vvanouytsel/jetspotter:dev

