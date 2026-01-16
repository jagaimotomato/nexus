#!/bin/bash

# 检查是否输入了项目名称
if [ -z "$1" ]; then
  echo "Usage: ./create-go-project.sh <project-name>"
  exit 1
fi

PROJECT_NAME=$1
MODULE_NAME=$1 # 你可以改为 github.com/yourname/$1

echo "🚀 正在按照 golang-standards/project-layout 初始化项目: $PROJECT_NAME ..."

# 1. 创建项目根目录
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME" || exit

# 2. 初始化 Go Module
go mod init "$MODULE_NAME"

# 3. 创建标准目录结构
# 核心应用代码
mkdir -p cmd/server
mkdir -p internal/biz      # 业务逻辑 (Business Logic)
mkdir -p internal/data     # 数据访问 (Data Access)
mkdir -p internal/service  # 接口实现 (Service Layer)
mkdir -p internal/conf     # 配置结构体定义

# 库代码
mkdir -p pkg/util          # 公共工具库 (可被外部引用)

# API 与 协议
mkdir -p api/protobuf      # gRPC proto 文件
mkdir -p api/swagger       # Swagger 文档

# 配置与构建
mkdir -p configs           # 配置文件 (yaml/json)
mkdir -p build/package     # Dockerfile 存放处
mkdir -p scripts           # 脚本 (Build/Deploy)
mkdir -p deployments       # K8s manifests, docker-compose

# Web 前端 (如果是一体化项目)
mkdir -p web/dist

# 文档与测试
mkdir -p docs
mkdir -p test

# 4. 创建基础文件
# .gitignore
cat > .gitignore <<EOF
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Environment variables
.env

# IDE specific files
.idea/
.vscode/
*.swp
EOF

# main.go
cat > cmd/server/main.go <<EOF
package main

import "fmt"

func main() {
	fmt.Println("Hello, $PROJECT_NAME! Project initialized based on golang-standards.")
}
EOF

# config.yaml (示例配置)
cat > configs/config.yaml <<EOF
server:
  port: 8080
  name: "$PROJECT_NAME"
database:
  driver: mysql
  source: root:123456@tcp(127.0.0.1:3306)/dbname
EOF

# Dockerfile (基础模板)
cat > build/package/Dockerfile <<EOF
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
CMD ["./server"]
EOF

# README.md
cat > README.md <<EOF
# $PROJECT_NAME

Standard Go project layout based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

## Directory Structure

- **cmd/**: Main applications for this project.
- **internal/**: Private application and library code.
- **pkg/**: Library code that's ok to use by external applications.
- **api/**: OpenAPI/Swagger specs, JSON schema files, protocol definition files.
EOF

echo "✅ 项目 $PROJECT_NAME 创建成功！"
echo "📂 进入目录: cd $PROJECT_NAME"
echo "🏃 运行项目: go run cmd/server/main.go"