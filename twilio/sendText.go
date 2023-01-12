package twilio

import (
	"fmt"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type Body struct {
	Body string `json:"body"`
	To   string `json:"to"`
}

func SendText(c *gin.Context) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))
	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	client := twilio.NewRestClient()

	org := c.MustGet("orgs").(model.Organization)

	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)

	var body Body
	c.BindJSON(&body)

	params := &api.CreateMessageParams{}
	params.SetBody(body.Body)
	params.SetFrom(org.TwilioNumber)
	params.SetTo(body.To)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}

	c.JSON(200, gin.H{
		"message": body.Body,
		"to":      body.To,
	})
}
