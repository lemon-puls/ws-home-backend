package main

import (
	"fmt"
	"ws-home-backend/config"
)

func main() {
	// 初始化配置
	config.InitConfig("./config/config-dev.yaml")
	fmt.Printf("Conf: %+v\n", config.Conf.MysqlConfig)
}
