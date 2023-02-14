package org

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Org struct {
	model.Organization
	Pos model.Pos
}

func GetOrg(c *gin.Context) {
	org, _ := c.Get("orgs")

	orgId := c.Param("orgId")

	orgIdInt, orgIdIntErr := strconv.Atoi(orgId)

	if orgIdIntErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid URL parameter",
		})
		return
	}

	var pos []model.Pos
	// get the pos that is associated with the org
	err := config.DB.Where(&model.Pos{OrganizationID: int64(orgIdInt)}).Find(&pos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if org == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No organization found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"org": org,
		"pos": pos,
	})
}
