package mock_backend

import (
	"github.com/gin-gonic/gin"
)

func NewMockServer() *gin.Engine {
	r := gin.New()

	// Mock endpoints
	r.GET("/api/events/1", func(c *gin.Context) {
		c.JSON(200, "")
	})

	r.GET("/api/echo-realip", func(c *gin.Context) {
		c.String(200, c.Request.Header.Get("X-Real-IP"))
	})

	r.GET("/api/echo-xff", func(c *gin.Context) {
		c.String(200, c.Request.Header.Get("X-Forwarded-For"))
	})

	return r
}
