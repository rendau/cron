package cmd

import (
	"github.com/rendau/cron/internals/domain/types"
	"github.com/rendau/dop/dopTools"
	"github.com/spf13/viper"
)

var conf = struct {
	Debug    bool           `mapstructure:"DEBUG"`
	LogLevel string         `mapstructure:"LOG_LEVEL"`
	ConfPath string         `mapstructure:"CONF_PATH"`
	Jobs     []*types.JobSt `mapstructure:"JOBS"`
}{}

func confLoad() {
	dopTools.SetViperDefaultsFromObj(conf)

	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	_ = viper.Unmarshal(&conf)
}
