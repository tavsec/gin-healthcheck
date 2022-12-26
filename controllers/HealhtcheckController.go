package controllers

import "github.com/gin-gonic/gin"

func HealthcheckController(c *gin.Context) {
	c.Status(200)
}
