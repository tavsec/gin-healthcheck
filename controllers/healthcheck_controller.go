package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

func HealthcheckController(checks []checks.Check, config config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		statuses := make([]CheckStatus, 0)
		httpStatus := config.StatusOK
		for _, check := range checks {
			pass := check.Pass()
			statuses = append(statuses, CheckStatus{
				Name: check.Name(),
				Pass: pass,
			})

			if !pass {
				httpStatus = config.StatusNotOK
			}
		}

		c.JSON(httpStatus, statuses)
	}

	return gin.HandlerFunc(fn)
}
