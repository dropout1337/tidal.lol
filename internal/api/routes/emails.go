package routes

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"tidal.lol/internal/emails"
	"tidal.lol/internal/smtpd"
)

// GetEmails returns a list of emails based on the provided email address or regex pattern or wildcard.
func GetEmails(c *gin.Context) {
	var emailList []smtpd.Email
	address := c.Param("email")

	// If the address starts with a wildcard, return all emails that match the domain.
	if strings.HasPrefix(address, "*") {
		for _, email := range emails.Emails {
			for _, to := range email.To {
				if strings.Split(to, "@")[1] == strings.Split(address, "@")[1] {
					emailList = append(emailList, email)
				}
			}
		}
	} else if r, err := regexp.Compile(address); err == nil { // If the address is a regex pattern, return all emails that match the pattern.
		for _, email := range emails.Emails {
			for _, to := range email.To {
				if r.MatchString(to) {
					emailList = append(emailList, email)
				}
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
