test: ## Run all the tests
	go test ./...

bench: ## Run all the tests
	go test -bench=. -benchmem

torture-test: ## Run all the tests
	export GOVALUATE_TORTURE_TEST="true"
	go test -bench=. -benchmem

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test torture-test help
.DEFAULT_GOAL := help