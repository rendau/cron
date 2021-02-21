package entities

type JobSt struct {
	Time          string `mapstructure:"time"`
	Url           string `mapstructure:"url"`
	RetryCount    int    `mapstructure:"retry_count"`
	RetryInterval int    `mapstructure:"retry_interval"`
}
