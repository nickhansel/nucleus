package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetPurchases(c *gin.Context) {
	purchases := []model.Purchase{}
	// add the purchased items and the variation related to the purchased item to the purchase struct
	err := config.DB.Preload("PurchasedItems").Preload("PurchasedItems.Variation").Preload("Location").Find(&purchases).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
	})
}
