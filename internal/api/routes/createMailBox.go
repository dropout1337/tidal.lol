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
		utils.HTTPResponse(c, 400)
		return
	}

	if len(data.Password) > 16 {
		utils.HTTPResponse(c, 400, "password too long")
		return
	}

	if len(data.Email) > 64 {
		utils.HTTPResponse(c, 400, "email too long")
		return
	}

	query, err := database.DB.Query("SELECT * FROM mailbox WHERE email=$email", map[string]any{"email": data.Email})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}
	result := utils.ToJson(query)

	if len(result.Get("0.result").Array()) != 0 {
		utils.HTTPResponse(c, 302, "mailbox already exists")
		return
	}

	mailbox, err := database.DB.Create("mailbox", map[string]interface{}{
		"email":    data.Email,
		"password": data.Password,
		"token":    uuid.NewString(),
	})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}

	utils.HTTPResponse(c, 200, map[string]string{
		"token": utils.ToJson(mailbox).Get("0.token").String(),
	})
}
