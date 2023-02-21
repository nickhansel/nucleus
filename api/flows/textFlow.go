package flows

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/cron/text"
	"github.com/nickhansel/nucleus/model"
	"strconv"
)

type TextFlowBody struct {
	To       []string `json:"to"`
	TextBody string   `json:"text_body"`
	Date     string   `json:"date"`
}

func ScheduleTextFlows(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var textFlowBody TextFlowBody
	if err := c.ShouldBindJSON(&textFlowBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	convertedIds := make([]int64, len(textFlowBody.To))
	for i, id := range textFlowBody.To {
		convertedIds[i], _ = strconv.ParseInt(id, 10, 64)
	}

	ids := convertedIds
	textBody := textFlowBody.TextBody
	date := textFlowBody.Date

	text.ScheduleFlowTexts(date, ids, org, textBody)

	c.JSON(200, gin.H{"message": "Texts scheduled"})

}
