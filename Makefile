SOURCES := $(shell find . -name '*.go')
BINARY := kubectl-eks

build: kubectl-eks

clean:
	@rm -rf $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags="-s -w" main.go
