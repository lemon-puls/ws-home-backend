# 定义输出目录
OUTPUT_DIR := output
# 获取当前 Git 提交的哈希值，格式为 8 位
GIT_HASH := $(shell git show -s --format=%h --abbrev=8)
# 定义二进制文件名称
BIN_NAME := ws-home

# 定义生成的工件路径，包括输出目录、二进制名称和 Git 哈希
ARTIFACT := ${OUTPUT_DIR}/${BIN_NAME}_${GIT_HASH}

# 声明伪目标，表示这些目标不对应实际文件
.PHONY: linux build gen docker build-docker

# 默认编译选项（windows）
build:                     # build 目标的命令
 	# 使用 Go 编译器构建项目，输出为 .exe 文件
	go build -o ${ARTIFACT}.exe

# 定义 linux 目标，设置环境变量以便于交叉编译
# 禁用 CGO
linux: export CGO_ENABLED=0
# 设置目标操作系统为 Linux
linux: export GOOS=linux
# 设置目标架构为 amd64
linux: export GOARCH=amd64
# linux 目标的命令
linux:
	# 使用 Go 编译器构建项目，输出到指定的工件路径
	go build ${BUILD_FLAGS} -o ${ARTIFACT}

docker: gen build-docker

gen:
	which swag || go install github.com/swaggo/swag/cmd/swag@latest
	go generate ./...


build-docker:
	go build ${BUILD_FLAGS} -o ${BIN_NAME} main.go


