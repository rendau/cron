package core

import (
	"github.com/rendau/cron/internals/domain/entities"
	"github.com/rendau/cron/internals/interfaces"
	"github.com/robfig/cron"
)

type St struct {
	lg   interfaces.Logger
	jobs []*entities.JobSt

	cron *cron.Cron
}

func New(lg interfaces.Logger, jobs []*entities.JobSt) *St {
	return &St{
		lg:   lg,
		jobs: jobs,

		cron: cron.New(),
	}
}
