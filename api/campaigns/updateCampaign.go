package campaign

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"regexp"
	"strconv"
	"time"
)

type UpdateCampaignBody struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func UpdateCampaign(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	campaignId := c.Query("campaignId")
	//convert to int64
	id, err := strconv.ParseInt(campaignId, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid campaignId"})
		return
	}

	var body UpdateCampaignBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var campaign model.Campaign
	config.DB.Where("\"id\" = ? AND \"organizationId\" = ?", id, org.ID).First(&campaign)

	if campaign.ID == 0 {
		c.JSON(400, gin.H{"error": "Campaign not found"})
		return
	}

	if body.Name != "" {
		campaign.Name = body.Name
	}
	if body.Status != "" {
		campaign.IsDeleted = body.Status == "DELETED"
	}

	config.DB.Save(&campaign)

	c.JSON(200, gin.H{"message": "Campaign updated"})
}

type SmsCampaignBody struct {
	SendTime string `json:"send_time"`
	Body     string `json:"body"`
}

func UpdateSMSCampaign(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	campaignId := c.Query("campaignId")
	//convert to int64
	id, err := strconv.ParseInt(campaignId, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid campaignId"})
		return
	}

	var body SmsCampaignBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var campaign model.TextCampaign
	config.DB.Where("\"id\" = ?", id).First(&campaign)

	var rootCampaign model.Campaign
	config.DB.Where("\"id\" = ?", campaign.CampaignID).First(&rootCampaign)
	if rootCampaign.OrganizationID != org.ID {
		c.JSON(400, gin.H{"error": "Campaign not found"})
		return
	}

	if campaign.ID == 0 {
		c.JSON(400, gin.H{"error": "Campaign not found"})
		return
	}

	if body.SendTime != "" {
		// validate sendTime is in the form 2023-02-13 10:46:20
		if ValidateTime(body.SendTime) == false {
			c.JSON(400, gin.H{"error": "Invalid sendTime"})
			return
		}
		//validate that the date is in the future
		if IsFutureDate(body.SendTime) == false {
			c.JSON(400, gin.H{"error": "Send time must be in the future"})
			return
		}

		campaign.SendTime = body.SendTime
	}
	if body.Body != "" {
		campaign.Body = body.Body
	}

	config.DB.Save(&campaign)

	c.JSON(200, gin.H{"message": "Campaign updated"})
}

type EmailCampaignBody struct {
	SendTime string `json:"send_time"`
	Subject  string `json:"subject"`
	HTML     string `json:"html"`
}

func UpdateEmailCampaign(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	campaignId := c.Query("campaignId")
	//convert to int64
	id, err := strconv.ParseInt(campaignId, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid campaignId"})
		return
	}

	var body EmailCampaignBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var campaign model.EmailCampaign
	config.DB.Where("\"id\" = ?", id).First(&campaign)

	var rootCampaign model.Campaign
	config.DB.Where("\"id\" = ?", campaign.CampaignID).First(&rootCampaign)
	if rootCampaign.OrganizationID != org.ID {
		c.JSON(400, gin.H{"error": "Campaign not found"})
		return
	}

	if campaign.ID == 0 {
		c.JSON(400, gin.H{"error": "Campaign not found"})
		return
	}

	if body.SendTime != "" {
		// validate sendTime is in the form 2023-02-13 10:46:20
		if ValidateTime(body.SendTime) == false {
			c.JSON(400, gin.H{"error": "Invalid sendTime"})
			return
		}
		//validate that the date is in the future
		if IsFutureDate(body.SendTime) == false {
			c.JSON(400, gin.H{"error": "Send time must be in the future"})
			return
		}

		campaign.SendTime = body.SendTime
	}
	if body.Subject != "" {
		campaign.Subject = body.Subject
		campaign.Text = body.Subject
	}
	if body.HTML != "" {
		campaign.HTML = body.HTML
	}

	c.JSON(200, gin.H{"message": "Campaign updated"})

}

func ValidateTime(time string) bool {
	//	validate sendTime is in the form 2023-02-13 10:46:20
	var valid = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	return valid.MatchString(time)
}

func IsFutureDate(dateStr string) bool {
	layout := "2006-01-02 15:04:05"

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Invalid date format")
		return false
	}

	now := time.Now()
	if t.After(now) {
		fmt.Println("Date is in the future")
		return true
	} else {
		fmt.Println("Date is in the past")
		return false
	}
}
