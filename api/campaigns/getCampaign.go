package campaign

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"strconv"
)

func GetCampaign(c *gin.Context) {
	id := c.Query("id")
	//convert id to int
	convertedString, _ := strconv.Atoi(id)
	org := c.MustGet("orgs").(model.Organization)

	//	check to make sure the campaign belongs to the org
	doesOrgBelongToCampaign := config.DB.Where("\"organizationId\" = ? AND \"id\" = ?", org.ID, convertedString).First(&model.Campaign{})

	if doesOrgBelongToCampaign.Error != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var result map[string]interface{}

	var campaign model.Campaign
	//find the campaign by id and dont include relations where id = 0
	config.DB.Preload("FbCampaign").Where("\"id\" = ?", convertedString).First(&campaign)

	if campaign.IsTextCampaign {
		var textCampaign model.TextCampaign
		config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&textCampaign)
		result = map[string]interface{}{
			"campaign":      campaign,
			"text_campaign": textCampaign,
		}
		//	add the text campaign to the campaign object
	} else {
		var emailCampaign model.EmailCampaign
		config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&emailCampaign)
		result = map[string]interface{}{
			"campaign":       campaign,
			"email_campaign": emailCampaign,
		}
	}

	var fbCampaign model.FbCampaign
	config.DB.Where("\"campaignId\" = ?", campaign.ID).First(&fbCampaign)
	if fbCampaign.ID != 0 {
		result = map[string]interface{}{
			"campaign":    campaign,
			"fb_campaign": fbCampaign,
		}
	}

	c.JSON(200, gin.H{"campaign": result})
}

func GetCampaignsById(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	campaignId := c.Param("campaignId")

	//conivrt id to int64
	convertedId, _ := strconv.ParseInt(campaignId, 10, 64)

	var campaigns model.Campaign
	config.DB.Where("\"organizationId\" = ? AND \"id\" = ?", org.ID, convertedId).First(&campaigns)

	if campaigns.ID == 0 {
		c.JSON(401, gin.H{"error": "No campaigns found"})
		return
	}

	if campaigns.IsTextCampaign {
		var textCampaign model.TextCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&textCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": textCampaign, "purchases": purchases})
		return
	}

	if campaigns.IsFbCampaign {
		var fbCampaign model.FbCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&fbCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": fbCampaign, "purchases": purchases})
	}

	if campaigns.IsEmailCampaign {
		var emailCampaign model.EmailCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&emailCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": emailCampaign, "purchases": purchases})
	}

}
