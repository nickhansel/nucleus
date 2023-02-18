package flows

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/cron/email"
	"github.com/nickhansel/nucleus/model"
)

type EmailFlowBody struct {
	To        []int64 `json:"to"`
	EmailBody string  `json:"email_body"`
	Subject   string  `json:"subject"`
	Date      string  `json:"date"`
}

func ScheduleEmailFlows(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var emailFlowBody EmailFlowBody
	if err := c.ShouldBindJSON(&emailFlowBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ids := emailFlowBody.To
	emailBody := emailFlowBody.EmailBody
	date := emailFlowBody.Date
	subject := emailFlowBody.Subject

	fmt.Println("ids: ", ids)

	email.ScheduleFlowEmails(date, ids, org, emailBody, subject)

	c.JSON(200, gin.H{"message": "Emails scheduled"})

}
