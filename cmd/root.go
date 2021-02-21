package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rendau/cron/internals/adapters/logger/zap"
	"github.com/rendau/cron/internals/domain/core"
	"github.com/rendau/cron/internals/domain/entities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "cron",
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		lg, err := zap.New(viper.GetString("log_level"), viper.GetBool("debug"), false)
		if err != nil {
			log.Fatal(err)
		}

		jobs := viper.Get("parsed_jobs").([]*entities.JobSt)

		core := core.New(lg, jobs)

		lg.Infow("Starting", "http_listen", viper.GetString("http_listen"))

		for _, job := range jobs {
			lg.Infow(
				"Cron job",
				"time", job.Time,
				"url", job.Url,
			)
		}

		err = core.StartCron()
		if err != nil {
			log.Fatal(err)
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		lg.Infow("Shutting down...")

		core.StopCron()

		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	viper.SetDefault("debug", "false")
	viper.SetDefault("http_listen", ":9090")
	viper.SetDefault("log_level", "debug")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	jobs := make([]*entities.JobSt, 0)
	err := viper.UnmarshalKey("jobs", &jobs)
	if err != nil {
		log.Fatal(err)
	}

	viper.Set("parsed_jobs", jobs)
}
