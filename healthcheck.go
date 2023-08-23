package gin_healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"github.com/tavsec/gin-healthcheck/controllers"
)

func New(engine *gin.Engine, config config.Config, checks []checks.Check) error {
	engine.Handle(config.Method, config.HealthPath, controllers.HealthcheckController(checks, config))
	return nil
}
