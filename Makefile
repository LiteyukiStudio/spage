PROJECT_NAME := spage
BUILD_DIR := build
REV := $(shell git rev-parse --short HEAD)
GO_BUILD := go build
GO_FLAGS := -trimpath
LD_FLAGS := -ldflags="-s -w"

# Go 环境变量
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 1

# 输出文件名
OUTPUT_NAME := $(BUILD_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH)

# 根据目标平台设置 Clang 交叉编译参数
ifeq ($(GOOS),darwin)
  ifeq ($(GOARCH),amd64)
    CC_TARGET := x86_64-apple-darwin
  else ifeq ($(GOARCH),arm64)
    CC_TARGET := aarch64-apple-darwin
  endif
  
  # 检查是否在 macOS 上构建
  UNAME_S := $(shell uname -s)
  ifeq ($(UNAME_S),Darwin)
    # 在 macOS 上直接使用系统 Clang
    CC := clang
    CLANG_FLAGS :=
  else
    # 在非 macOS 上使用交叉编译
    ifeq ($(shell which xcrun 2>/dev/null),)
      # 不是 macOS 环境，尝试使用 osxcross (如果已安装)
      ifneq ($(shell which o64-clang 2>/dev/null),)
        CC := o64-clang
        CLANG_FLAGS :=
      else
        # 降级到不使用 CGO
        CGO_ENABLED := 0
        CC := clang
        CLANG_FLAGS :=
      endif
    else
      # macOS 环境，使用 xcrun
      CC := xcrun clang
      CLANG_FLAGS :=
    endif
  endif
else ifeq ($(GOOS),windows)
  ifeq ($(GOARCH),amd64)
    CC_TARGET := x86_64-w64-mingw32
  else ifeq ($(GOARCH),386)
    CC_TARGET := i686-w64-mingw32
  endif
  
  # 检查是否有 MinGW 工具链
  ifneq ($(shell which $(CC_TARGET)-gcc 2>/dev/null),)
    CC := $(CC_TARGET)-gcc
    CLANG_FLAGS :=
  else ifneq ($(shell which clang 2>/dev/null),)
    CC := clang
    CLANG_FLAGS := --target=$(CC_TARGET)
  else
    # 降级到不使用 CGO
    CGO_ENABLED := 0
    CC := gcc
    CLANG_FLAGS :=
  endif
else ifeq ($(GOOS),linux)
  ifeq ($(GOARCH),amd64)
    CC_TARGET := x86_64-linux-gnu
  else ifeq ($(GOARCH),386)
    CC_TARGET := i386-linux-gnu
  else ifeq ($(GOARCH),arm64)
    CC_TARGET := aarch64-linux-gnu
  else ifeq ($(GOARCH),arm)
    CC_TARGET := arm-linux-gnueabihf
  endif
  
  # 检查特定 Linux 交叉编译工具链
  ifneq ($(shell which $(CC_TARGET)-gcc 2>/dev/null),)
    CC := $(CC_TARGET)-gcc
    CLANG_FLAGS :=
  else ifneq ($(shell which clang 2>/dev/null),)
    CC := clang
    CLANG_FLAGS := --target=$(CC_TARGET)
  else
    # 使用系统默认编译器
    CC := gcc
    CLANG_FLAGS :=
  endif
else
  # 默认使用系统 Clang
  CC := clang
  CLANG_FLAGS :=
endif

# 主构建目标
.PHONY: all clean spage

all: spage

spage:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(PROJECT_NAME)-$(GOOS)-$(GOARCH):$(REV) for $(GOOS)/$(GOARCH)"
	@if [ "$(CGO_ENABLED)" = "1" ]; then \
		CC="$(CC) $(CLANG_FLAGS)" \
		CGO_ENABLED=1 \
		GOOS=$(GOOS) \
		GOARCH=$(GOARCH) \
		$(GO_BUILD) $(GO_FLAGS) $(LD_FLAGS) -o $(OUTPUT_NAME) ./cmd/server; \
	else \
		echo "Note: Building without CGO (CGO_ENABLED=0)"; \
		CGO_ENABLED=0 \
		GOOS=$(GOOS) \
		GOARCH=$(GOARCH) \
		$(GO_BUILD) $(GO_FLAGS) $(LD_FLAGS) -o $(OUTPUT_NAME) ./cmd/server; \
	fi

.PHONY: web
web:
	cd web-src && pnpm install && pnpm build
	mkdir -p ./$(BIN_NAME)/static/dist
	cp -r web-src/out/* ./$(BIN_NAME)/static/dist

.PHONY: proto
proto:
	protoc --go_out=protos/result --go_opt=paths=source_relative --go-grpc_out=protos/result --go-grpc_opt=paths=source_relative protos/source/*.proto