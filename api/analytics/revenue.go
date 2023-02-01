package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetTotalRevenue(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var purchases []model.Purchase
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&purchases)

	var campaigns []model.Campaign
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&campaigns)

	attributedRevenue := 0.0
	for _, campaign := range campaigns {
		attributedRevenue += campaign.AttributedRevenue
	}

	totalRevenue := 0.0
	for _, purchase := range purchases {
		totalRevenue += purchase.AmountMoney
	}

	c.JSON(200, gin.H{
		"total_revenue":      totalRevenue,
		"attributed_revenue": attributedRevenue,
	})

}
