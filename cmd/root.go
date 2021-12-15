package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rendau/cron/internals/adapters/logger/zap"
	"github.com/rendau/cron/internals/domain/core"
	"github.com/rendau/cron/internals/domain/entities"
	"github.com/rendau/cron/internals/interfaces"
	"github.com/spf13/viper"
)

func Execute() {
	var err error

	app := struct {
		lg   interfaces.Logger
		core *core.St
	}{}

	loadConfig()

	app.lg, err = zap.New(viper.GetString("LOG_LEVEL"), viper.GetBool("DEBUG"), false)
	if err != nil {
		log.Fatal(err)
	}

	jobs := viper.Get("PARSED_JOBS").([]*entities.JobSt)

	app.core = core.New(app.lg, jobs)

	app.lg.Infow("Starting", "http_listen", viper.GetString("HTTP_LISTEN"))

	for _, job := range jobs {
		app.lg.Infow(
			"Cron job",
			"time", job.Time,
			"url", job.Url,
		)
	}

	err = app.core.StartCron()
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.lg.Infow("Shutting down...")

	app.core.StopCron()

	os.Exit(0)
}

func loadConfig() {
	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("HTTP_LISTEN", ":9090")
	viper.SetDefault("LOG_LEVEL", "debug")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	jobs := make([]*entities.JobSt, 0)
	err := viper.UnmarshalKey("JOBS", &jobs)
	if err != nil {
		log.Fatal(err)
	}

	viper.Set("PARSED_JOBS", jobs)
}
