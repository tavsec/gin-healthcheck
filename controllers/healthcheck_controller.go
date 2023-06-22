package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"golang.org/x/sync/errgroup"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

var errHealthcheckFailed = errors.New("healthcheck failed")

func HealthcheckController(checks []checks.Check, config config.Config) gin.HandlerFunc {
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
					return errHealthcheckFailed
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			httpStatus = config.StatusNotOK
		}

		c.JSON(httpStatus, statuses)
	}

	return gin.HandlerFunc(fn)
}
