package routes

import (
	"github.com/gin-gonic/gin"
	"tidal.lol/internal/database"
	"tidal.lol/internal/utils"
)

// GetEmails returns all emails for a mailbox.
func GetEmails(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.HTTPResponse(c, 401)
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE token=$authorizationToken", map[string]any{"authorizationToken": token})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}
	result := utils.ToJson(query)

	if len(result.Get("0.result").Array()) == 0 {
		utils.HTTPResponse(c, 400, "mailbox not found")
		return
	}

	emails, err := database.DB.Query("SELECT * FROM inbox WHERE token=$authorizationToken", map[string]any{"authorizationToken": token})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}

	for i := 0; i < len(emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{})); i++ {
		delete(emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{})[i].(map[string]interface{}), "token")
	}

	utils.HTTPResponse(c, 200, emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{}))
}
