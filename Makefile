bin/ojosama: go.* *.go cmd/*
	make test
	mkdir -p bin
	go fmt .
	go build -o bin/ojosama ./cmd/ojosama

.PHONY: test
test:
	go test -cover ./...
