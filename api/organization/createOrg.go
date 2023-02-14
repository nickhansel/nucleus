package org

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func CreateNewPos(posName string, orgId int64, c *gin.Context) (model.Pos, error) {
	// get the date right now

	pos := model.Pos{
		Name:           posName,
		OrganizationID: orgId,
		AccessToken:    "",
		RefreshToken:   "",
		ExpiresAt:      time.Now(),
		MerchantID:     "",
	}

	err := config.DB.Create(&pos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	return pos, nil
}

func CreateOrg(c *gin.Context) {
	// get the user id from the previous middleware
	id, _ := c.Get("id")

	// check if the user already has an organization
	var user model.User
	err := config.DB.Where("id = ?", id).First(&user).Error

	fmt.Println(id, "user")

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already has an organization",
		})
		return
	}

	// get the name froom the request body
	type Body struct {
		Name    string `json:"name" binding:"required"`
		OrgType string `json:"type" binding:"required"`
		Pos     string `json:"pos" binding:"required"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	name := body.Name
	orgType := body.OrgType
	pos := body.Pos

	// create a new organization
	org := model.Organization{
		Name: name,
		Type: orgType,
		// FbCustomAudiences is an array of strings
		TwilioNumber:     "",
		SendgridEmail:    "",
		IsSendgridAuthed: false,
		IsTwilioAuthed:   false,
		PosID:            0,
		// add the user as a member of the org
		Members: []model.User{
			{
				ID: id.(int64),
			},
		},
	}

	// create the org
	createOrgErr := config.DB.Create(&org).Error

	if createOrgErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create a new POS
	newPos, err := CreateNewPos(pos, org.ID, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// update the org with the new pos id
	org.PosID = newPos.ID

	err = config.DB.Save(&org).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"org":     org,
		"pos":     newPos,
	})

}
