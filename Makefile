BIN_NAME ?= spage
GO_PKG_ROOT ?= github.com/LiteyukiStudio/spage
GO_ENTRYPOINT_SERVER ?= ./cmd/server
GO_ENTRYPOINT_AGENT ?= ./cmd/agent

GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
ZIGARCH=$(shell \
    arch="$(GOARCH)"; \
    if [ "$$arch" = "amd64" ]; then echo "x86_64"; \
    elif [ "$$arch" = "arm64" ]; then echo "aarch64"; \
    elif [ "$$arch" = "386" ]; then echo "i386"; \
    elif [ "$$arch" = "arm" ]; then echo "arm"; \
    else echo "$$arch"; fi \
)

ZIGOS=$(shell \
    os="$(GOOS)"; \
    if [ "$$os" = "windows" ]; then echo "windows"; \
    elif [ "$$os" = "darwin" ]; then echo "macos"; \
    elif [ "$$os" = "linux" ]; then echo "linux"; \
    elif [ "$$os" = "freebsd" ]; then echo "freebsd"; \
    else echo "$$os"; fi \
)

ZIGABI=$(shell \
    if [ "$(GOOS)" = "linux" ]; then \
        if [ "$(GOARCH)" = "arm" ]; then echo "gnueabihf"; \
        else echo "gnu"; fi; \
    elif [ "$(GOOS)" = "windows" ]; then echo "gnu"; \
    else echo "none"; fi \
)

ZIGTARGET=$(ZIGARCH)-$(ZIGOS)-$(ZIGABI)

.PHONY: web
web:
	cd web-src && pnpm install && pnpm build
	mkdir -p ./$(BIN_NAME)/static/dist
	cp -r web-src/out/* ./$(BIN_NAME)/static/dist

.PHONY: proto
proto:
	protoc --go_out=protos/result --go_opt=paths=source_relative --go-grpc_out=protos/result --go-grpc_opt=paths=source_relative protos/source/*.proto

.PHONY: spage
spage:
	@mkdir -p build
	@( \
	OUTNAME=$(BIN_NAME)-$(GOOS)-$(GOARCH); \
	VERSION=$$(git describe --tags --always 2>/dev/null || echo dev); \
	echo "Building $$OUTNAME:$$VERSION for $(GOOS)/$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then OUTNAME=$${OUTNAME}.exe; fi; \
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	CC="zig cc -target $(ZIGTARGET)" CXX="zig c++ -target $(ZIGTARGET)" \
	go build -trimpath \
	-ldflags "-X '$(GO_PKG_ROOT)/config.CommitHash=$$(git rev-parse HEAD)' \
	-X '$(GO_PKG_ROOT)/config.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
	-X '$(GO_PKG_ROOT)/config.Version=$${VERSION}'" \
	-o build/$${OUTNAME} $(GO_ENTRYPOINT_SERVER) \
	)