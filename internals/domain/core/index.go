package core

import (
	"github.com/rendau/cron/internals/domain/types"
	"github.com/rendau/dop/adapters/logger"
	"github.com/robfig/cron"
)

type St struct {
	lg   logger.Lite
	jobs []*types.JobSt

	cron *cron.Cron
}

func New(lg logger.Lite, jobs []*types.JobSt) *St {
	return &St{
		lg:   lg,
		jobs: jobs,

		cron: cron.New(),
	}
}
