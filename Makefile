.PHONY := all fmt vet

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

.DEFAULT_GOAL :=
build: vet
	go build
