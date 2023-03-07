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

	var customerGroup model.CustomerGroup
	config.DB.Where("\"id\" = ?", campaigns.CustomerGroupID).First(&customerGroup)

	if campaigns.ID == 0 {
		c.JSON(401, gin.H{"error": "No campaigns found"})
		return
	}

	if campaigns.IsTextCampaign {
		var textCampaign model.TextCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&textCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": textCampaign, "purchases": purchases, "customer_group": customerGroup})
		return
	}

	if campaigns.IsFbCampaign {
		var fbCampaign model.FbCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&fbCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": fbCampaign, "purchases": purchases, "customer_group": customerGroup})
	}

	if campaigns.IsEmailCampaign {
		var emailCampaign model.EmailCampaign
		config.DB.Where("\"campaignId\" = ?", campaigns.ID).First(&emailCampaign)

		var purchases []model.Purchase
		config.DB.Where("\"attributedCampaignId\" = ?", campaigns.ID).Find(&purchases)

		var emailMetrics []model.EmailCampaignAnalytics
		config.DB.Where("\"emailCampaignId\" = ?", emailCampaign.ID).Find(&emailMetrics)

		type EmailMetricData struct {
			TotalSent         int32 `json:"total_sent"`
			TotalDelivered    int32 `json:"total_delivered"`
			TotalBounces      int32 `json:"total_bounces"`
			TotalClicks       int32 `json:"total_clicks"`
			TotalUniqueClicks int32 `json:"total_unique_clicks"`
			TotalOpens        int32 `json:"total_opens"`
			TotalUniqueOpens  int32 `json:"total_unique_opens"`
			TotalSpamReports  int32 `json:"total_spam_reports"`
			TotalBlocked      int32 `json:"total_blocked"`
			TotalUnsubscribes int32 `json:"total_unsubscribes"`
			TotalInvalid      int32 `json:"total_invalid"`
		}

		data := EmailMetricData{}

		if len(emailMetrics) > 0 {

			for _, metric := range emailMetrics {
				data.TotalSent += metric.Sent
				data.TotalDelivered += metric.Delivered
				data.TotalBounces += metric.Bounces
				data.TotalClicks += metric.Clicks
				data.TotalUniqueClicks += metric.UniqueClicks
				data.TotalOpens += metric.Opens
				data.TotalUniqueOpens += metric.UniqueOpens
				data.TotalSpamReports += metric.SpamReports
				data.TotalBlocked += metric.Blocked
				data.TotalUnsubscribes += metric.Unsubscribed
				data.TotalInvalid += metric.Invalid
			}
		}
		c.JSON(200, gin.H{"campaign": campaigns, "sub_campaign": emailCampaign, "purchases": purchases, "metrics": data, "customer_group": customerGroup})
	}

}
