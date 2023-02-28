package items

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func GetAllItems(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var items []model.Item
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&items)

	c.JSON(200, gin.H{
		"code":  "SUCCESS",
		"items": items,
	})
}
