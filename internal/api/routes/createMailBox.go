package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tidal.lol/internal/database"
	"tidal.lol/internal/utils"
)

type CreateMailBoxRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateMailBox creates a new mailbox.
func CreateMailBox(c *gin.Context) {
	var data CreateMailBoxRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.DefaultResponse(c, 400)
		return
	}

	if len(data.Password) > 16 {
		utils.DefaultResponse(c, 400, "password too long")
		return
	}

	if len(data.Email) > 64 {
		utils.DefaultResponse(c, 400, "email too long")
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE email=$email", map[string]any{"email": data.Email})
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	if len(query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})) != 0 {
		utils.DefaultResponse(c, 302, "mailbox already exists")
		return
	}

	mailbox, err := database.DB.Create("mailbox", map[string]interface{}{
		"email":    data.Email,
		"password": data.Password,
		"token":    uuid.NewString(),
	})
	if err != nil {
		utils.DefaultResponse(c, 500)
		return
	}

	delete(mailbox.([]interface{})[0].(map[string]interface{}), "password")
	utils.DefaultResponse(c, 200, mailbox)
}
