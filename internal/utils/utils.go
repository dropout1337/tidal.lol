package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"net/http"
)

func HTTPResponse(c *gin.Context, status int, data ...any) {
	if len(data) > 0 {
		c.AbortWithStatusJSON(200, gin.H{"code": status, "message": data[0]})
	} else {
		c.AbortWithStatusJSON(200, gin.H{"code": status, "message": fmt.Sprintf("%v: %v", status, http.StatusText(status))})
	}
}

func ToJson(v interface{}) gjson.Result {
	query, _ := json.Marshal(v)
	return gjson.Parse(string(query))
}
