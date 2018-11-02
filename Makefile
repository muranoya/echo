GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=gofmt
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY=echo

.PHONY: test build clean check fmt
all: fmt check test build

run:
	start_server --port 8080 --pid-file=./echo.pid -- ./echo

fmt:
	$(GOFMT) -s -l -e -w .

check:
	errcheck -exclude errcheck_excludes.txt -asserts -verbose ./...
	go vet ./...
	golint

test:
	$(GOTEST) -v ./...

build:
	$(GOBUILD) -o $(BINARY) cmd/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY)
