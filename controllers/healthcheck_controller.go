package controllers

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"golang.org/x/sync/errgroup"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

var ErrHealthcheckFailed = errors.New("healthcheck failed")

func HealthcheckController(checks []checks.Check, config config.Config) gin.HandlerFunc {
	var lock sync.Mutex
	var failureInARow uint32

	fn := func(c *gin.Context) {
		var eg errgroup.Group

		statuses := make([]CheckStatus, len(checks))
		httpStatus := config.StatusOK
		for idx, check := range checks {
			captureCheck := check
			captureIdx := idx
			eg.Go(func() error {
				pass := captureCheck.Pass()
				statuses[captureIdx] = CheckStatus{
					Name: captureCheck.Name(),
					Pass: pass,
				}

				if !pass {
					return ErrHealthcheckFailed
				}
				return nil
			})
		}

		lock.Lock()
		if err := eg.Wait(); err != nil {
			httpStatus = config.StatusNotOK
			failureInARow += 1

			if failureInARow >= config.FailureNotification.Threshold &&
				config.FailureNotification.Chan != nil {
				config.FailureNotification.Chan <- err
			}
		} else {
			if failureInARow != 0 && config.FailureNotification.Chan != nil {
				failureInARow = 0
				config.FailureNotification.Chan <- nil
			}
		}
		lock.Unlock()

		c.JSON(httpStatus, statuses)
	}

	return gin.HandlerFunc(fn)
}
