package gin_healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/rsmnarts/gin-healthcheck/checks"
	"github.com/rsmnarts/gin-healthcheck/config"
	"github.com/rsmnarts/gin-healthcheck/controllers"
)

func New(engine *gin.Engine, config config.Config, checks []checks.Check) error {
	engine.Handle(config.Method, config.HealthPath, controllers.HealthcheckController(checks, config))
	return nil
}
