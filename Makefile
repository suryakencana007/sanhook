PKG := sanhook

build:
	@go build -o ${PKG} $(shell PWD)/cmd/${PKG}/main.go

run:
	@go run $(shell PWD)/cmd/${PKG}/main.go http

test:
	@go test -cover ./...

.PHONY: run test
