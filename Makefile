BIN_NAME ?= spage
GO_PKG_ROOT ?= github.com/LiteyukiStudio/spage
GO_ENTRYPOINT_SERVER ?= ./cmd/server
GO_ENTRYPOINT_AGENT ?= ./cmd/agent

GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)

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
	# zig cc 
	ZIGARCH=$(GOARCH) \
	if [ "$(GOARCH)" = "amd64" ]; then ZIGARCH=x86_64; fi; \
	if [ "$(GOARCH)" = "arm64" ]; then ZIGARCH=aarch64; fi; \
	if [ "$(GOARCH)" = "386" ]; then ZIGARCH=i386; fi; \
	if [ "$(GOARCH)" = "arm" ]; then ZIGARCH=arm; fi; \

	ZIGOS=$(GOOS) \
	if [ "$(GOOS)" = "windows" ]; then ZIGOS=windows; fi; \
	if [ "$(GOOS)" = "darwin" ]; then ZIGOS=macos; fi; \
	if [ "$(GOOS)" = "linux" ]; then ZIGOS=linux; fi; \
	if [ "$(GOOS)" = "freebsd" ]; then ZIGOS=freebsd; fi; \
	
	ZIGTARGET=$${ZIGOS}-$${ZIGARCH} \

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