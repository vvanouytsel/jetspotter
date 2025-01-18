config_snippet := "docs/snippets/config.snippet"
output_snippet := "docs/snippets/output.snippet"
overview_snippet := "docs/snippets/overview.snippet"

just:
    just --list

build: 
	go build -o jetspotter -ldflags "-linkmode external -extldflags -static" cmd/jetspotter/jetspotter.go

run:
	go run cmd/jetspotter/*
	
unit-tests:
	go test $(go list ./... | grep -v postgres)

integration-tests:
	go test $(go list ./... | grep postgres)

generate-doc:
	echo "✨ Generating documentation..."
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

forward-database: apply-manifests
	echo "✨ You can connect to your local database via: kubectl exec -ti -n dev $(kubectl get pods  -l app=postgres --no-headers -o custom-columns=":metadata.name" -n dev) -- psql -U jetspotter"
	echo "✨ Alternatively you can use 'psql -U jetspotter -h localhost -p 5432'"
	kubectl port-forward -n dev service/postgres 5432:5432

apply-manifests:
	just check kubectl
	kubectl apply -f development/
	kubectl get pods --no-headers -o custom-columns=":metadata.name" -n dev | xargs -I {} kubectl wait --for=condition=Ready pod/{} -n dev 

load-image: build-container
	minikube image load ghcr.io/vvanouytsel/jetspotter:dev

create-dev:
	#!/bin/bash
	set -e
	just check minikube
	minikube start --driver docker
	just load-image
	just apply-manifests
	just forward-database

destroy-dev:
	minikube delete

check $tool:
	#!/bin/bash
	set -e
	if ! command -v $tool > /dev/null; then
		echo "✨ Please install $tool."
		exit 1
	fi
