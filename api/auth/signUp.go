package auth

import (
	"github.com/gin-gonic/gin"
	token "github.com/nickhansel/nucleus/api/utils/auth"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Body struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OrgName   string `json:"org_name"`
}

func SignUp(c *gin.Context) {
	body := Body{}
	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	hashPassword, err := token.HashPassword(body.Password)
	if err != nil {
		return
	}

	org := model.Organization{
		Name: body.OrgName,
		Type: "AGENCY",
		// FbCustomAudiences is an array of strings
		TwilioNumber:     "",
		SendgridEmail:    "",
		IsSendgridAuthed: false,
		IsTwilioAuthed:   false,
		PosID:            0,
		// add the user as a member of the or
	}

	config.DB.Create(&org)

	user := model.User{
		Email:          body.Email,
		Password:       hashPassword,
		FirstName:      body.FirstName,
		LastName:       body.LastName,
		OrganizationID: org.ID,
	}

	config.DB.Create(&user)

	//update the member of the org
	err = config.DB.Model(&org).Association("Members").Append(&user)
	if err != nil {
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"org":     org,
		"user":    user,
	})

}
