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
			var customerGroup model.CustomerGroup
			config.DB.Where("\"id\" = ?", campaign.CustomerGroupID).First(&customerGroup)

			result = append(result, map[string]interface{}{
				"campaign":      campaign,
				"text_campaign": textCampaign,
				"group_name":    customerGroup.Name,
				"group_id":      customerGroup.ID,
			})
		} else if campaign.IsEmailCampaign {
			var emailCampaign model.EmailCampaign
			config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&emailCampaign)
			var customerGroup model.CustomerGroup
			config.DB.Where("\"id\" = ?", campaign.CustomerGroupID).First(&customerGroup)

			result = append(result, map[string]interface{}{
				"campaign":       campaign,
				"email_campaign": emailCampaign,
				"group_name":     customerGroup.Name,
				"group_id":       customerGroup.ID,
			})
		} else if campaign.IsFbCampaign {
			var fbCampaign model.FbCampaign
			config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&fbCampaign)
			var customerGroup model.CustomerGroup
			config.DB.Where("\"id\" = ?", campaign.CustomerGroupID).First(&customerGroup)

			result = append(result, map[string]interface{}{
				"campaign":    campaign,
				"fb_campaign": fbCampaign,
				"group_name":  customerGroup.Name,
				"group_id":    customerGroup.ID,
			})
		}
	}

	c.JSON(200, gin.H{
		"code":      "SUCCESS",
		"campaigns": result,
	})
}
