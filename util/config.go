package util

import (
	"time"

	"github.com/spf13/viper"
)

// 保存环境变量
type Config struct {
	DBDriver           string        `mapstructure:"DB_DRIVER"`
	DBSource           string        `mapstructure:"DB_SOURCE"`
	ServerAddress      string        `mapstructure:"SERVER_ADDRESS"`
	TokenKey           string        `mapstructure:"TOKEN_KEY"`
	AcessTokenDuration time.Duration `mapstructure:"ACESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv() //更新配置文件用的

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	// fmt.Println(config)
	return
}
