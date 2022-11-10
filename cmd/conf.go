package cmd

import (
	"os"

	"github.com/rendau/cron/internals/domain/types"
	"github.com/rendau/dop/dopTools"
	"github.com/spf13/viper"
)

var conf = struct {
	Debug    bool           `mapstructure:"DEBUG"`
	LogLevel string         `mapstructure:"LOG_LEVEL"`
	Jobs     []*types.JobSt `mapstructure:"JOBS"`
}{}

func confLoad() {
	dopTools.SetViperDefaultsFromObj(conf)

	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("LOG_LEVEL", "info")

	confPath := os.Getenv("CONF_PATH")
	if confPath == "" {
		confPath = "conf.yml"
	}
	viper.SetConfigFile(confPath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	_ = viper.Unmarshal(&conf)
}
