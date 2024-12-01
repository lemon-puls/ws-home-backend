# WS HOME
## 项目简介

本项目是一个基于 Go 语言开发的相册管理系统后端服务。
> [前往前端项目](https://github.com/lemon-puls/ws-home-backend)

主要功能：

- 用户管理：支持用户注册、登录、信息更新等基础功能，使用 JWT 进行身份认证
- 相册管理：支持创建相册、上传照片和视频、删除媒体文件等功能
- 对象存储：使用腾讯云 COS 对象存储服务来存储照片和视频文件
- 数据统计：支持统计用户的相册数量、照片数量、存储空间使用情况等
- API 文档：集成 Swagger 文档，方便接口调试和查看

技术特点：

- 使用 Gin 框架构建 RESTful API
- 采用 GORM 作为 ORM 框架操作 MySQL 数据库
- 集成 Redis 用于缓存和令牌管理
- 使用雪花算法生成分布式 ID
- 支持日志分割和多环境配置
- 使用 Make 工具进行项目构建，支持跨平台编译

项目遵循良好的工程实践，包括：
- 统一的错误处理和响应格式
- 请求参数验证和国际化
- 中间件实现的日志记录和异常恢复
- 支持优雅关闭和重启

## Swagger
1. 执行以下命令生成 swagger 文档：
```bash
swag init
```
2. 启动后端服务，访问 `http://localhost:8080/swagger/index.html` 即可查看 swagger 文档。

> [官方使用说明文档](https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%A3%B0%E6%98%8E%E5%BC%8F%E6%B3%A8%E9%87%8A%E6%A0%BC%E5%BC%8F)

## 部署
1. 编译项目：

本项目使用了 make 进行项目的构建、编译，命令如下：
```bash
# 默认编译为 windows 平台下的可执行文件
make # 或者 make build
# 编译为 linux 平台下的可执行文件
make linux
```
> 安装 make:  
[官网下载地址](https://sourceforge.net/projects/gnuwin32/files/make/3.81/make-3.81-bin.zip/download?use_mirror=zenlayer&download=)

2. 把的可执行文件上传到服务器上
在同一目录下创建指向可执行文件的软链接，命名为 ws-home-backend, 因为 app.sh 会根据这个名称来启动项目。
```bash
# 创建软链接
ln -sf ws-home-backend_<替换为实际 git commit hash> ws-home-backend
```
3. 使用脚本启动项目
把 script/app.sh 上传到服务器上, 和可执行文件放在同一目录下
```bash
# 赋予脚本执行权限
chmod +x app.sh
# 启动项目
./app.sh start
# 其他操作
# 停止项目
./app.sh stop
# 重启项目
./app.sh restart
# 查看状态
./app.sh status
```