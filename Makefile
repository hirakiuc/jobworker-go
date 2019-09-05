OUTPUT := "jobdaemon"
GOPATH_BIN := $(GOPATH)/bin

.PHONY: deps
deps:
	go mod download

.PHONY: clean
clean:
	go clean
	rm -f $(OUTPUT)

check:
	golangci-lint run --enable-all -D dupl ./...

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o $(OUTPUT) ./cmd/jobdaemon/main.go

.PHONY: imagebuild
imagebuild:
	docker build . -t test

default: build
