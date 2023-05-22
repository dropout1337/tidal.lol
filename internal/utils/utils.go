package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"net/http"
)

func DefaultResponse(c *gin.Context, status int, data ...any) {
	if len(data) > 0 {
		c.AbortWithStatusJSON(200, gin.H{"code": status, "message": data[0]})
	} else {
		c.AbortWithStatusJSON(200, gin.H{"code": status, "message": fmt.Sprintf("%v: %v", status, http.StatusText(status))})
	}

}
