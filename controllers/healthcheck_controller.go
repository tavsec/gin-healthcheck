package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
)

type CheckStatus struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

func HealthcheckController(checks []checks.Check) gin.HandlerFunc {
	statuses := make([]CheckStatus, 0)
	for _, check := range checks {
		pass := check.Pass()
		statuses = append(statuses, CheckStatus{
			Name: check.Name(),
			Pass: pass,
		})
	}

	fn := func(c *gin.Context) {
		c.JSON(200, statuses)
	}

	return gin.HandlerFunc(fn)
}
