package config

import (
	"gitee.com/wxlao/config-client"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

const (
	JAEGER_SERVER  = "JAEGER_SERVER"
	SESSION_ID_KEY = "SESSION_ID_KEY"
	REQUEST_ID_KEY = "REQUEST_ID_KEY"
	USER_ID_KEY    = "USER_ID_KEY"
	APP_ID         = "APP_ID"
	ENV_TYPE       = "ENV_TYPE"
)

// Config 配置
type Config struct {
	viper *viper.Viper
}

// New 新配置
func New(viper *viper.Viper) *Config {
	return &Config{
		viper: viper,
	}
}

// Load 加载配置
func (c *Config) Load() error {
	c.viper.SetConfigType("env")

	if err := config.LoadFile(".env"); err != nil {
		log.Error().Str("err", err.Error()).Msg("Error reading config file")
	}

	c.setDefault()

	return nil
}

func (c *Config) setDefault() {
	c.viper.SetDefault(JAEGER_SERVER, "")
	c.viper.SetDefault(SESSION_ID_KEY, "Ayg-Sessionid")
	c.viper.SetDefault(REQUEST_ID_KEY, "logcontext-requestid")
	c.viper.SetDefault(USER_ID_KEY, "logcontext-userid")
	c.viper.SetDefault(APP_ID, "an-app")
	c.viper.SetDefault(ENV_TYPE, "pro")
}
