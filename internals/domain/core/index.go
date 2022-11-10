package core

import (
	"net/http"
	"time"

	"github.com/rendau/cron/internals/domain/types"
	"github.com/rendau/dop/adapters/logger"
	"github.com/robfig/cron/v3"
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

		cron: cron.New(
			cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month)),
		),
	}
}

func (c *St) Start() error {
	var err error

	for _, job := range c.jobs {
		if job.Method == "" {
			job.Method = "GET"
		}
		if job.Timeout <= 0 {
			job.Timeout = 5 * time.Second
		}
		if job.RetryCount < 0 {
			job.RetryCount = 0
		}

		jobInst := job // fix reference of job

		_, err = c.cron.AddFunc(job.Time,
			func() {
				c.handler(jobInst)
			},
		)
		if err != nil {
			c.lg.Errorw("Cron error", err)
			return err
		}
	}

	c.cron.Start()

	return nil
}

func (c *St) handler(job *types.JobSt) {
	defer func() {
		if err := recover(); err != nil {
			c.lg.Errorw("Recover", err)
		}
	}()

	for i := 0; i <= job.RetryCount; i++ {
		if c.sendReq(job) == nil {
			break
		}

		if job.RetryInterval > 0 {
			time.Sleep(job.RetryInterval)
		}
	}
}

func (c *St) sendReq(job *types.JobSt) error {
	req, err := http.NewRequest(job.Method, job.Url, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return nil
	}

	httpClient := http.Client{Timeout: job.Timeout}

	rep, err := httpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err, "url", job.Url)
		return err
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode, "url", job.Url)
		return nil
	}

	return nil
}

func (c *St) StopAndWait() {
	ctx := c.cron.Stop()
	<-ctx.Done()
}
