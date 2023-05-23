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
		utils.HTTPResponse(c, 400)
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE id=$id", map[string]any{"id": data.ID})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}
	result := utils.ToJson(query)

	if len(result.Get("0.result").Array()) == 0 {
		utils.HTTPResponse(c, 404, "mailbox not found")
		return
	} else if result.Get("0.result.0.password").String() != data.Password {
		utils.HTTPResponse(c, 401, "invalid password")
		return
	}

	_, err = database.DB.Delete(data.ID)
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}

	utils.HTTPResponse(c, 204)
}
