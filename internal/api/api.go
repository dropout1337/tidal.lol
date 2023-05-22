package api

import (
	"fmt"
	ginRecovery "github.com/LaysDragon/gin-custom-recovery-handler"
	"github.com/gin-gonic/gin"
	"tidal.lol/internal/api/routes"
	"tidal.lol/internal/logging"
	"time"
)

func NewServer(port int, debug bool) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.ForwardedByClientIP = true
	router.RemoveExtraSlash = true

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{"code": 404, "message": "404: Not Found"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(405, gin.H{"code": 405, "message": "405: Method not allowed"})
	})

	router.Use(ginRecovery.Recovery(func(c *gin.Context, serverErr interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"code": 500, "message": "500: Internal server error"})
	}))

	if debug {
		router.Use(func(c *gin.Context) {
			start := time.Now()
			c.Next()

			logging.Logger.Debug().
				Str("method", c.Request.Method).
				Str("uri", c.Request.RequestURI).
				Int("status", c.Writer.Status()).
				Str("ip", c.RemoteIP()).
				Str("time", time.Since(start).String()).
				Msg("Request received")
		})
	}

	v1 := router.Group("/api/v1")
	v1.Any("/emails/:email", routes.GetEmails)

	err := router.Run(fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}

	return nil
}
