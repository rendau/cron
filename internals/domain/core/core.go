package core

import (
	"errors"
	"net/http"
	"time"
)

func (c *St) StartCron() error {
	for _, v := range c.jobs {
		timer := v.Time
		url := v.Url
		retryCount := v.RetryCount
		retryInterval := v.RetryInterval

		err := c.cron.AddFunc(timer,
			func() {
				c.handler(url, retryCount, retryInterval)
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

func (c *St) StopCron() {
	c.cron.Stop()
}

func (c *St) handler(url string, retryCount int, retryInterval int) {
	defer func() {
		if r := recover(); r != nil {
			c.lg.Errorw("Recover", r)
		}
	}()

	for retryI := 0; retryI < retryCount; retryI++ {
		err := c.sendReq(url)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(retryInterval) * time.Second)
	}
}

func (c *St) sendReq(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return err
	}

	httpClient := http.Client{Timeout: 20 * time.Second}

	rep, err := httpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return err
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return errors.New("bad status code")
	}

	return nil
}
