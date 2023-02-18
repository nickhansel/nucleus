package flows

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"strconv"
)

const (
	ACTIVE = "ACTIVE"
	PAUSED = "PAUSED"
)

func UpdateFlowStatus(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	id := c.Param("flowId")
	// convert flowId to int64
	flowId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid flowId"})
		return
	}
	status := c.Query("status")

	if status != ACTIVE && status != PAUSED {
		c.JSON(400, gin.H{"error": "Invalid status"})
		return
	}

	var flow model.Flow
	if err := config.DB.Where("\"id\" = ? AND \"organizationId\" = ?", flowId, org.ID).First(&flow).Error; err != nil {
		c.JSON(400, gin.H{"error": "Flow not found"})
		return
	}

	flow.Status = status
	config.DB.Save(&flow)

	c.JSON(200, gin.H{"message": "Flow status updated"})
}
