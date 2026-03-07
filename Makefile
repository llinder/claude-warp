BINARY = claude-warp
PLUGIN_SCRIPTS = plugins/warp/scripts
GOFLAGS = -trimpath
LDFLAGS = -s -w

.PHONY: all build install clean test lint

all: build

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(PLUGIN_SCRIPTS)/$(BINARY) ./cmd/claude-warp

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(PLUGIN_SCRIPTS)/$(BINARY)

# Cross-compile for distribution
dist:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o dist/$(BINARY)-darwin-amd64 ./cmd/claude-warp
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o dist/$(BINARY)-darwin-arm64 ./cmd/claude-warp
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o dist/$(BINARY)-linux-amd64 ./cmd/claude-warp
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o dist/$(BINARY)-linux-arm64 ./cmd/claude-warp
