package api

import (
	"fmt"
	ginRecovery "github.com/LaysDragon/gin-custom-recovery-handler"
	"github.com/gin-gonic/gin"
	"tidal.lol/internal/api/routes"
	"tidal.lol/internal/logging"
	"tidal.lol/internal/utils"
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
		utils.DefaultResponse(c, 404)
	})

	router.NoMethod(func(c *gin.Context) {
		utils.DefaultResponse(c, 405)
	})

	router.Use(ginRecovery.Recovery(func(c *gin.Context, serverErr interface{}) {
		utils.DefaultResponse(c, 500)
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
	v1.Any("/emails/:email", routes.GetTempEmails)

	v2 := router.Group("/api/v2")

	v2.POST("/emails/create", routes.CreateMailBox)
	v2.DELETE("/emails/delete", routes.DeleteMailBox)
	v2.GET("/emails/inbox", routes.GetEmails)

	err := router.Run(fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}

	return nil
}
