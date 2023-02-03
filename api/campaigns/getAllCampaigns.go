package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetAllCampaigns(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var campaigns []model.Campaign
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&campaigns)

	var result []map[string]interface{}

	for _, campaign := range campaigns {
		if campaign.IsTextCampaign {
			var textCampaign model.TextCampaign
			config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&textCampaign)
			result = append(result, map[string]interface{}{
				"campaign":      campaign,
				"text_campaign": textCampaign,
			})
		} else if campaign.IsEmailCampaign {
			var emailCampaign model.EmailCampaign
			config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&emailCampaign)
			result = append(result, map[string]interface{}{
				"campaign":      campaign,
				"text_campaign": emailCampaign,
			})
		} else {
			var fbCampaign model.FbCampaign
			config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&fbCampaign)
			result = append(result, map[string]interface{}{
				"campaign":      campaign,
				"text_campaign": fbCampaign,
			})
		}
	}

	c.JSON(200, gin.H{
		"code":      "SUCCESS",
		"campaigns": result,
	})
}
