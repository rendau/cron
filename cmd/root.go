package cmd

import (
	"log"
	"os"

	"github.com/rendau/cron/internals/domain/core"
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	"github.com/rendau/dop/dopTools"
)

func Execute() {
	var err error

	app := struct {
		lg   *dopLoggerZap.St
		core *core.St
	}{}

	confLoad()

	app.lg = dopLoggerZap.New(conf.LogLevel, conf.Debug)

	app.core = core.New(
		app.lg,
		conf.Jobs,
	)

	for _, job := range conf.Jobs {
		app.lg.Infow(
			"Cron job",
			"time", job.Time,
			"url", job.Url,
		)
	}

	// START

	app.lg.Infow("Starting")

	err = app.core.Start()
	if err != nil {
		log.Fatal(err)
	}

	var exitCode int

	<-dopTools.StopSignal()

	// STOP

	app.lg.Infow("Shutting down...")

	app.lg.Infow("Wait routines...")

	// app.core.StopAndWait()

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}
