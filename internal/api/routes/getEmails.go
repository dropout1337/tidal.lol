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
		utils.DefaultResponse(c, 401)
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE token=$authorizationToken", map[string]any{"authorizationToken": token})
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	if len(query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})) == 0 {
		utils.DefaultResponse(c, 400, "mailbox not found")
		return
	}

	emails, err := database.DB.Query("SELECT * FROM inbox WHERE token=$authorizationToken", map[string]any{"authorizationToken": token})
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	for i := 0; i < len(emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{})); i++ {
		delete(emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{})[i].(map[string]interface{}), "token")
	}

	utils.DefaultResponse(c, 200, emails.([]interface{})[0].(map[string]interface{})["result"].([]interface{}))
}
