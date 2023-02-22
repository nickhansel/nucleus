package auth

import (
	"github.com/gin-gonic/gin"
	token "github.com/nickhansel/nucleus/api/utils/auth"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/nickhansel/nucleus/sendinblue"
	"strconv"
)

type InviteBody struct {
	Email string `json:"email"`
}

func SendInviteEmail(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var invite model.Invite
	var body InviteBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	invite.Email = body.Email
	invite.OrganizationID = org.ID

	if err := config.DB.Create(&invite).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sendinblue.SendInviteEmail(body.Email, org, invite.ID)

	c.JSON(200, gin.H{"message": "Invite sent"})

}

type AcceptInviteBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func AcceptInvite(c *gin.Context) {
	var invite model.Invite

	//convert param id to int64
	convertedID, err := strconv.ParseInt(c.Query("id"), 10, 64)

	if err := config.DB.Where("id = ?", convertedID).First(&invite).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var body AcceptInviteBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := token.HashPassword(body.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:          body.Email,
		Password:       hashPassword,
		FirstName:      body.FirstName,
		LastName:       body.LastName,
		OrganizationID: invite.OrganizationID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var org model.Organization
	if err := config.DB.Where("id = ?", invite.OrganizationID).First(&org).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = config.DB.Model(&org).Association("Members").Append(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	invite.IsAccepted = true
	config.DB.Save(&invite)

	c.JSON(200, gin.H{
		"message": "Invite accepted",
		"user":    user,
	})
}
