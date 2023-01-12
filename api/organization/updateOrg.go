package org

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	Name string `json:"name" binding:"required"`
}

func UpdateOrg(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	// get the name froom the request body
	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get the org that is associated with the user
	var organization model.Organization
	err := config.DB.Where(&model.Organization{ID: org.ID}).First(&organization).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// update the name of the org
	organization.Name = body.Name

	// save the changes
	err = config.DB.Save(&organization).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"org": organization,
	})

}
