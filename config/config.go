package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	*ServerConfig    `mapstructure:"server"`
	*LogConfig       `mapstructure:"log"`
	*MysqlConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*SnowflakeConfig `mapstructure:"snowflake"`
	*AuthConfig      `mapstructure:"auth"`
	*CosConfig       `mapstructure:"cos"`
	*AmapConfig      `mapstructure:"amap"`
	*EmailConfig     `mapstructure:"email"`
}

type ServerConfig struct {
	Name    string `mapstructure:"name"`
	Port    int    `mapstructure:"port"`
	Profile string `mapstructure:"profile"`
	Version string `mapstructure:"version"`
}

type CosConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
}

type AuthConfig struct {
	JwtExpire time.Duration `mapstructure:"jwt_expire"`
	JwtSecret string        `mapstructure:"jwt_secret"`
}

type SnowflakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type AmapConfig struct {
	Key      string `mapstructure:"key"`
	RegeoURL string `mapstructure:"regeo_url"`
}

type EmailConfig struct {
	FromEmail    string   `mapstructure:"from_email"`
	FromPassword string   `mapstructure:"from_password"`
	ToEmails     []string `mapstructure:"to_emails"`
	SmtpServer   string   `mapstructure:"smtp_server"`
	SmtpPort     int      `mapstructure:"smtp_port"`
}

func InitConfig(path string) {
	// 设置配置文件路径
	viper.SetConfigFile(path)
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		panic(err)
	}

	// 将读取的配置绑定到结构体变量
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		panic(err)
	}

	// 监听配置文件变化
	viper.WatchConfig()
	// 注册回调函数 当配置文件变化时 更新配置
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s", e.Name)
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
			panic(err)
		}
	})
}
