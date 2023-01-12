package sendgrid

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/sendgrid/sendgrid-go"
)

func ResendVerificationEmail(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var Organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&Organization)

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for sendgrid")
	}

	sendGridID := fmt.Sprintf("%v", Organization.SendgridID)

	apiKey := os.Getenv("SENDGRID_API_KEY")

	host := "https://api.sendgrid.com"

	request := sendgrid.GetRequest(apiKey, "/v3/verified_senders/resend/"+sendGridID, host)

	request.Method = "POST"

	response, err := sendgrid.API(request)

	statusCode := response.StatusCode

	if err != nil {
		log.Println(err)
	}

	if statusCode != 204 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "Error sending verification email",
			"response": statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Verification email sent",
		"response": statusCode,
	})
}
