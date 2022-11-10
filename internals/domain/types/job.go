package types

import (
	"time"
)

type JobSt struct {
	Time          string        `mapstructure:"time"`
	Method        string        `mapstructure:"method"`
	Url           string        `mapstructure:"url"`
	Timeout       time.Duration `mapstructure:"timeout"`
	RetryCount    int           `mapstructure:"retry_count"`
	RetryInterval time.Duration `mapstructure:"retry_interval"`
}
