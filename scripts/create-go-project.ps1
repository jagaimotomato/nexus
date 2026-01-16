param (
    [Parameter(Mandatory=$true)]
    [string]$ProjectName
)

$ModuleName = $ProjectName

# 使用纯英文提示，避免 PowerShell 编码乱码问题
Write-Host ">> Initializing Go project: $ProjectName ..." -ForegroundColor Cyan

# 1. 创建并进入目录
New-Item -Path $ProjectName -ItemType Directory -Force | Out-Null
Set-Location $ProjectName

# 2. 初始化 Go Mod
go mod init $ModuleName

# 3. 定义目录列表
$dirs = @(
    "cmd/server",
    "internal/biz",
    "internal/data",
    "internal/service",
    "internal/conf",
    "pkg/util",
    "api/protobuf",
    "api/swagger",
    "configs",
    "build/package",
    "scripts",
    "deployments",
    "web/dist",
    "docs",
    "test"
)

# 创建目录
foreach ($dir in $dirs) {
    New-Item -Path $dir -ItemType Directory -Force | Out-Null
}

# 4. 创建基础文件

# .gitignore
$gitignoreContent = @"
# Binaries
*.exe
*.dll
*.so
*.dylib
*.test
*.out
vendor/
.env
.idea/
.vscode/
"@
Set-Content -Path .gitignore -Value $gitignoreContent -Encoding UTF8

# main.go
$mainContent = @"
package main

import "fmt"

func main() {
	fmt.Println("Hello, $ProjectName! Project initialized based on golang-standards.")
}
"@
Set-Content -Path cmd/server/main.go -Value $mainContent -Encoding UTF8

# config.yaml
$configContent = @"
server:
  port: 8080
  name: "$ProjectName"
"@
Set-Content -Path configs/config.yaml -Value $configContent -Encoding UTF8

# Dockerfile
$dockerContent = @"
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
"@
Set-Content -Path build/package/Dockerfile -Value $dockerContent -Encoding UTF8

# 结束提示 (纯英文)
Write-Host "SUCCESS: Project $ProjectName created!" -ForegroundColor Green
Write-Host "-> CD into project: cd $ProjectName"
Write-Host "-> Run app: go run cmd/server/main.go"