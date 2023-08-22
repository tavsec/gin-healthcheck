package controllers

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rsmnarts/gin-healthcheck/checks"
	"github.com/rsmnarts/gin-healthcheck/config"
	"golang.org/x/sync/errgroup"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

var (
	ErrHealthcheckFailed = errors.New("healthcheck failed")

	lock          sync.Mutex
	failureInARow uint32
)

func Healthcheck(checks []checks.Check, config config.Config) (int, []CheckStatus) {
	var (
		eg errgroup.Group

		statuses   = make([]CheckStatus, len(checks))
		httpStatus = config.StatusOK
	)

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

	return httpStatus, statuses
}

func HealthcheckController(checks []checks.Check, config config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(Healthcheck(checks, config))
	}

	return gin.HandlerFunc(fn)
}
