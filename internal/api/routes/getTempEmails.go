package routes

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"tidal.lol/internal/database"
	"tidal.lol/internal/emails"
	"tidal.lol/internal/smtpd"
	"tidal.lol/internal/utils"
)

// GetTempEmails returns a list of emails based on the provided email address or regex pattern or wildcard.
func GetTempEmails(c *gin.Context) {
	var emailList []smtpd.Email
	address := c.Param("email")

	// Check if the email address has been registered using the v2 api.
	query, err := database.DB.Query("SELECT * FROM mailbox WHERE email=$email", map[string]any{"email": address})
	if err != nil {
		utils.HTTPResponse(c, 500)
		return
	}
	result := utils.ToJson(query)

	if len(result.Get("0.result").Array()) != 0 {
		utils.HTTPResponse(c, 403, "this email address has been registered using the v2 api")
		return
	}

	// If the address starts with a wildcard, return all emails that match the domain.
	if strings.HasPrefix(address, "*") {
		for _, email := range emails.Emails {
			if strings.Split(email.To, "@")[1] == strings.Split(address, "@")[1] {
				emailList = append(emailList, email)
			}
		}
	} else if r, err := regexp.Compile(address); err == nil { // If the address is a regex pattern, return all emails that match the pattern.
		for _, email := range emails.Emails {
			if r.MatchString(email.To) {
				emailList = append(emailList, email)
			}
		}
	} else if regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`).MatchString(address) { // If the address is a valid email address, return all emails that match the address.
		emailList = emails.Get(address)
	}

	c.AbortWithStatusJSON(200, gin.H{
		"code":   200,
		"emails": emailList,
	})
}
