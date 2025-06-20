BIN_NAME ?= spage
GO_PKG_ROOT ?= github.com/LiteyukiStudio/spage
GO_ENTRYPOINT_SERVER ?= ./cmd/server
GO_ENTRYPOINT_AGENT ?= ./cmd/agent

GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
GOARM   ?=
GOAMD64 ?=
GO386   ?=
output  ?= $(GOARCH)

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
	OUTNAME=$(BIN_NAME)-$(GOOS)-$(output); \
	VERSION=$$(git describe --tags --always 2>/dev/null || echo dev); \
	echo "Building $$OUTNAME:$$VERSION for $(GOOS)/$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then OUTNAME=$${OUTNAME}.exe; fi; \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) GOAMD64=$(GOAMD64) GO386=$(GO386) \
	go build -trimpath \
	-ldflags "-X '$(GO_PKG_ROOT)/config.CommitHash=$$(git rev-parse HEAD)' \
	-X '$(GO_PKG_ROOT)/config.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
	-X '$(GO_PKG_ROOT)/config.Version=$${VERSION}'" \
	-o build/$${OUTNAME} $(GO_ENTRYPOINT_SERVER) \
	)

.PHONY: spage-container
spage-container: web spage
	@echo "Building container image for $(GOOS)/$(GOARCH)"; \
	docker build -t $(BIN_NAME):$(GOOS)-$(output) --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) .

# Apple native Linux container build, macOS 26 only
.PHONY: spage-apple-container
spage-apple-container: web spage
	@echo "Building Apple container image for $(GOOS)/$(GOARCH)"; \
	container build -t $(BIN_NAME):$(GOOS)-$(output) --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH)

.PHONY: agent
agent:
	@echo "Building agent for $(GOOS)/$(GOARCH)";

.PHONY: plugin

# 用法: make plugin name=XXX
plugin:
	@mkdir -p build
	@ext=so; \
	if [ "$(GOOS)" = "darwin" ]; then ext=dylib; fi; \
	echo "Building plugin: plugins/$(name) -> build/$(name).$${ext}"; \
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -buildmode=plugin \
	-o build/$(name).$${ext} ./plugins/$(name)