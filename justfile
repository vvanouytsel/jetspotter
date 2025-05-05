config_snippet := "docs/snippets/config.snippet"
output_snippet := "docs/snippets/output.snippet"
overview_snippet := "docs/snippets/overview.snippet"

just:
    just --list

build: 
	go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go

run:
	MAX_RANGE_KILOMETERS=70 go run cmd/jetspotter/*

run-bru:
	MAX_RANGE_KILOMETERS=100  LOCATION_LATITUDE="50.898706483327565" LOCATION_LONGITUDE="4.487719500638235" go run cmd/jetspotter/*
	
run-close:
	MAX_RANGE_KILOMETERS=20 LOCATION_LATITUDE="50.898706483327565" LOCATION_LONGITUDE="4.487719500638235" go run cmd/jetspotter/*
	
run-with-extended-scan:
	MAX_RANGE_KILOMETERS=30 MAX_SCAN_RANGE_KILOMETERS=40 FETCH_INTERVAL=10 go run cmd/jetspotter/*

test:
	go test ./...

generate-doc:
	echo "✨ Generating documentation..."
	go doc -u configuration.Config | sed -n '/type Config struct {/,/}/p' > {{config_snippet}}
	cat {{config_snippet}}
	go doc -u jetspotter.Aircraft | sed -n '/type Aircraft struct {/,/}/p' > {{output_snippet}}
	cat {{output_snippet}}
	echo -e "# Overview\n" > helm/jetspotter/README.md
	echo -e "[![GitHub repository](https://img.shields.io/badge/GitHub-jetspotter-green?logo=github)](https://github.com/vvanouytsel/jetspotter)\n" >> helm/jetspotter/README.md
	cat {{overview_snippet}} >> helm/jetspotter/README.md
	cat helm/jetspotter/README.md

build-container:
	#!/bin/bash
	set -e
	ct=$(command -v docker || command -v podman)
	[[ "$ct" == "" ]] && echo "✨ Please install Podman or Docker." && exit 1
	$ct build -t ghcr.io/vvanouytsel/jetspotter:dev .

run-container: build-container
	#!/bin/bash
	set -e
	ct=$(command -v docker || command -v podman)
	[[ "$ct" == "" ]] && echo "✨ Please install Podman or Docker." && exit 1
	$ct run ghcr.io/vvanouytsel/jetspotter:dev

run-web-container: build-container
	#!/bin/bash
	set -e
	ct=$(command -v docker || command -v podman)
	[[ "$ct" == "" ]] && echo "✨ Please install Podman or Docker." && exit 1
	echo "✨ Running jetspotter container with web UI on http://localhost:8080"
	$ct run -p 8080:8080 -p 8085:8085 -e WEB_UI_ENABLED=true ghcr.io/vvanouytsel/jetspotter:dev

apply-manifests:
	just check kubectl
	kubectl apply -f development/
	kubectl get pods --no-headers -o custom-columns=":metadata.name" -n dev | xargs -I {}  kubectl wait --for=condition=Ready pod/{} -n dev 

load-image: build-container
	minikube image load ghcr.io/vvanouytsel/jetspotter:dev

create-dev:
	#!/bin/bash
	set -e
	just check minikube
	minikube start --driver docker
	just load-image
	just apply-manifests

	echo "✨ You can connect to your local database via: kubectl exec -ti -n dev $(kubectl get pods  -l app=postgres --no-headers -o custom-columns=":metadata.name" -n dev) -- psql -U postgres"

destroy-dev:
	minikube delete

check $tool:
	#!/bin/bash
	set -e
	if ! command -v $tool > /dev/null; then
		echo "✨ Please install $tool."
		exit 1
	fi
