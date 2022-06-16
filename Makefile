bin/ojosama: test
	mkdir -p bin
	go fmt .
	go build -o bin/ojosama ./cmd/ojosama

.PHONY: test
test: go.* *.go cmd/*
	go test -cover ./...
