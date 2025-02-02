PLUGIN_DIR := plugins
BIN_DIR := bin/plugins
GOBUILD := go build -ldflags="-s -w" -o

# 获取所有插件目录
PLUGIN_DIRS := $(shell find $(PLUGIN_DIR) -type d -mindepth 1 -maxdepth 1)
# 从目录路径中提取插件名称
PLUGINS := $(notdir $(PLUGIN_DIRS))
# 生成目标文件路径
TARGETS := $(addprefix $(BIN_DIR)/,$(PLUGINS))

# 根据操作系统设置可执行文件后缀
ifeq ($(OS),Windows_NT)
    EXT := .exe
else
    EXT :=
endif

.PHONY: all plugins clean install help proto

all: plugins

plugins: $(TARGETS)

# 编译规则
$(BIN_DIR)/%: $(PLUGIN_DIR)/%/main.go
	@echo "Building plugin $*"
	@mkdir -p $(BIN_DIR)
	@cd $(PLUGIN_DIR)/$* && $(GOBUILD) ../../$@$(EXT) main.go

clean:
	@echo "Cleaning binaries..."
	@rm -rf $(BIN_DIR)

install: proto
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available targets:"
	@echo "  all       - Build all plugins (default)"
	@echo "  plugins   - Build all plugins"
	@echo "  clean     - Remove built binaries"
	@echo "  install   - Install dependencies"
	@echo "  help      - Show this help message"

# 调试用的目标，显示变量值
debug:
	@echo "Plugin directories found: $(PLUGIN_DIRS)"
	@echo "Plugin names: $(PLUGINS)"
	@echo "Target files: $(TARGETS)"

proto:
	@echo "Generating protobuf code..."
	@buf generate
