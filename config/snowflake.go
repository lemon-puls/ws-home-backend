package config

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"time"
)

var snowflakeNode *snowflake.Node

func InitSnowflakeNode(config *SnowflakeConfig) {

	machineID := config.MachineID
	startTime := config.StartTime

	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.L().Error("Error parsing start time", zap.Error(err))
		panic(err)
	}

	// 设置节点时间
	snowflake.Epoch = st.UnixNano() / 1000000

	snowflakeNode, err = snowflake.NewNode(machineID)

	if err != nil {
		zap.L().Error("Error creating snowflake", zap.Error(err))
		panic(err)
	}

	zap.L().Info("Snowflake node initialized")
}

// GenerateID 生成ID
func GenerateID() int64 {
	return snowflakeNode.Generate().Int64()
}
