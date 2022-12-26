package gin_healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/controllers"
)

func New(engine *gin.Engine, config Config) error {
	engine.GET(config.HealthPath, controllers.HealthcheckController)
	return nil
}
