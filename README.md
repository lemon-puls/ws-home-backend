# WS HOME

## Swagger
1. 执行以下命令生成 swagger 文档：
```bash
swag init
```
2. 启动后端服务，访问 `http://localhost:8080/swagger/index.html` 即可查看 swagger 文档。

> 使用文档：https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%A3%B0%E6%98%8E%E5%BC%8F%E6%B3%A8%E9%87%8A%E6%A0%BC%E5%BC%8F

## 部署
1. 编译后端项目：
```bash
# 编译在 Linux 下运行的可执行文件
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o output/ws-home-backend
```