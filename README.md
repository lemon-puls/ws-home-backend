# WS HOME

## Swagger
1. 执行以下命令生成 swagger 文档：
```bash
swag init
```
2. 启动后端服务，访问 `http://localhost:8080/swagger/index.html` 即可查看 swagger 文档。

> 使用文档：https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%A3%B0%E6%98%8E%E5%BC%8F%E6%B3%A8%E9%87%8A%E6%A0%BC%E5%BC%8F

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