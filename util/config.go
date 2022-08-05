package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Mode          string        `mapstructure:"GIN_MODE"`
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	TokenKey      string        `mapstructure:"TOKEN_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	fmt.Println(config)
	return
}
