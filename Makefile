OUTPUT := "jobdaemon"
GOPATH_BIN := $(GOPATH)/bin

.PHONY: dep
dep:
	dep ensure

.PHONY: clean
clean:
	go clean
	rm -f $(OUTPUT)

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o $(OUTPUT)

default: build
