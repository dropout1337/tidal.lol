package routes

import (
	"github.com/gin-gonic/gin"
	"tidal.lol/internal/database"
	"tidal.lol/internal/utils"
)

type DeleteMailBoxRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// DeleteMailBox deletes a mailbox.
func DeleteMailBox(c *gin.Context) {
	var data DeleteMailBoxRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.DefaultResponse(c, 400)
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE id=$id", map[string]any{"id": data.ID})
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	if len(query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})) == 0 {
		utils.DefaultResponse(c, 404, "mailbox not found")
		return
	} else if query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["password"] != data.Password {
		utils.DefaultResponse(c, 401, "invalid password")
		return
	}

	_, err = database.DB.Delete(data.ID)
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	utils.DefaultResponse(c, 204)
}
