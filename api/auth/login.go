package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	auth "github.com/nickhansel/nucleus/api/utils/auth"
	utils "github.com/nickhansel/nucleus/api/utils/token"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(c *gin.Context) {
	// get the body of the request
	var login Login
	err := c.ShouldBindJSON(&login)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get the user from the db
	var user model.User
	err = config.DB.Where("email = ?", login.Email).First(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	// check if the password is correct
	if !auth.CheckPasswordHash(login.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	// generate the jwt token
	accessToken, accessErr := utils.GenerateAccessToken(user.ID)
	refreshToken, refreshErr := utils.GenerateRefreshToken(user.ID)

	if accessErr != nil || refreshErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating token",
		})
		return
	}

	// send the token to the client
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}
