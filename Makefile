DOC_TEMPLATE = docs/snippets/config.snippet
.PHONY: run documentation

test:
	@echo "Running Golang tests..."
	go test ./...

doc:
	@echo "Generating documentation..."
	go doc -u configuration.Config | sed -n '/type Config struct {/,/}/p' > ${DOC_TEMPLATE}
	@cat ${DOC_TEMPLATE}
