OUTPUT := "jobdaemon"
GOPATH_BIN := $(GOPATH)/bin

dep:
	dep ensure

build: clean
	go build -o $(OUTPUT)

clean:
	go clean
	rm -f $(OUTPUT)

default: build
