bin/ojosama: go.* *.go cmd/* internal/*
	make test
	go vet .
	go fmt .
	mkdir -p bin
	go build -o bin/ojosama ./cmd/ojosama

.PHONY: test
test:
	go test -cover ./...

.PHONY: install
install: go.* *.go cmd/* internal/*
	go install ./cmd/ojosama
