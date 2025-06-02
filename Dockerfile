# Go 构建阶段
FROM golang:1.21 as backend-builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 设置 GOPROXY
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖项
RUN go mod download

# 复制后端源代码
COPY . .

# 构建应用程序
RUN make docker && \
    chmod +x ws-home && \
    ls -hail ws-home && \
    touch build_complete

# 最终运行阶段
FROM debian:bookworm-slim as final

# 安装CA证书和时区数据 否则最终容器无法通过 https 访问外部接口
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata && \
    rm -rf /var/lib/apt/lists/*

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata

# 更新CA证书
RUN update-ca-certificates

# 设置工作目录
WORKDIR /app

# 复制构建产物
COPY --from=backend-builder /app/build_complete /app/build_complete
COPY --from=backend-builder /app/ws-home /app/ws-home

# 清理构建标记文件
RUN rm build_complete && \
    ls -hail ws-home

# 暴露端口
EXPOSE 8080

# 启动命令
ENTRYPOINT ["./ws-home"]