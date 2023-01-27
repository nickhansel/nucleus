package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetEmailAnalytics(c *gin.Context) {
	campaignId := c.Param("email_campaign_id")

	// get the email campaign analytics
	var emailCampaigns []model.EmailCampaign
	config.DB.Preload("EmailCampaignAnalytics").Where("id = ?", campaignId).Find(&emailCampaigns)

	c.JSON(200, gin.H{
		"code":      "SUCCESS",
		"analytics": emailCampaigns,
	})
}
