# gcoord-go Makefile

.PHONY: build install clean test help

# 默认目标
all: build

# 构建可执行文件
build:
	@echo "🔨 构建 gcoord-go..."
	go build -o ./bin/gcoord ./cmd/gcoord
	@echo "✅ 构建完成: ./gcoord"

# 安装到系统
install:
	@echo "📦 安装 gcoord-go..."
	go install ./cmd/gcoord
	@echo "✅ 安装完成"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -f gcoord
	@echo "✅ 清理完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test ./...

# 运行基准测试
bench:
	@echo "⚡ 运行基准测试..."
	go test -bench=. ./...

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...

# 检查代码
vet:
	@echo "🔍 检查代码..."
	go vet ./...

# 下载依赖
deps:
	@echo "📥 下载依赖..."
	go mod tidy
	go mod download

# 显示帮助
help:
	@echo "gcoord-go 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  build   构建可执行文件"
	@echo "  install 安装到系统"
	@echo "  clean   清理构建文件"
	@echo "  test    运行测试"
	@echo "  bench   运行基准测试"
	@echo "  fmt     格式化代码"
	@echo "  vet     检查代码"
	@echo "  deps    下载依赖"
	@echo "  help    显示此帮助信息"
