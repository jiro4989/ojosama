bin/ojosama: go.* *.go cmd/* internal/*
	make test
	go vet .
	go fmt .
	mkdir -p bin
	go build -o bin/ojosama ./cmd/ojosama

.PHONY: test
test:
	go test -cover ./...
	which gocyclo && ./scripts/test_cyclomatic_complexity.sh

.PHONY: install
install: go.* *.go cmd/* internal/*
	go install ./cmd/ojosama

.PHONY: setup-tools
setup-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
