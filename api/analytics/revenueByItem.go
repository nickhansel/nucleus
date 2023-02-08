package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"strconv"
)

func GetReveneuByItem(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	itemId := c.Query("item_id")

	if itemId == "" {
		c.JSON(500, gin.H{
			"code":    "ERROR",
			"message": "Item id is required",
		})
		return
	}

	//convert item id to int
	ItemIdInt, err := strconv.Atoi(itemId)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    "ERROR",
			"message": "Error converting item id to int",
		})
		return
	}

	var purchases []model.Purchase
	err = config.DB.Where("\"organizationId\" = ?", org.ID).Preload("PurchasedItems").Preload("PurchasedItems.Variation").Preload("Customer").Find(&purchases).Error

	if err != nil {
		c.JSON(500, gin.H{
			"code":    "ERROR",
			"message": "Error getting purchases",
		})
		return
	}

	total := 0.0
	item := model.PurchasedItem{}

	//	find purchases where
	for _, purchase := range purchases {
		for _, purchasedItem := range purchase.PurchasedItems {
			if purchasedItem.VariationID == int64(ItemIdInt) {
				total += purchasedItem.Cost * float64(purchasedItem.Quantity)
				item = purchasedItem
			}
		}
	}

	if item.ID == 0 {
		c.JSON(200, gin.H{
			"code":    "SUCCESS",
			"message": "No purchases found",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":          "SUCCESS",
		"message":       "Found purchases",
		"total_revenue": total,
		"item":          item,
	})
	return

}
