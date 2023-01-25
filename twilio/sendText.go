package twilio

import (
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
	Body string  `json:"body"`
	To   []int32 `json:"to"`
}

func SendTextAPI(c *gin.Context) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))

	client := twilio.NewRestClient()

	org := c.MustGet("orgs").(model.Organization)

	var organization model.Organization
	config.DB.Where("id = ?", org.ID).First(&organization)

	numbers := []string{}

	// get the numbers associated with the customer ids in the to field of the body

	var body Body
	c.BindJSON(&body)

	for _, number := range body.To {
		var customer model.Customer
		config.DB.Where("id = ?", number).First(&customer)
		numbers = append(numbers, customer.PhoneNumber)
	}

	for index := range body.To {
		params := &api.CreateMessageParams{}
		params.SetBody(body.Body)
		params.SetFrom(organization.TwilioNumber)
		params.SetTo(numbers[index])

		resp, err := client.Api.CreateMessage(params)

		if err != nil || resp.ErrorCode != nil {
			log.Fatal(err)
		}

		c.JSON(200, gin.H{
			"message": "success",
		})

	}
}

func SendScheduledTexts(TextCampaign model.TextCampaign, org model.Organization) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))

	client := twilio.NewRestClient()

	organization := org

	numbers := TextCampaign.TargetNumbers
	body := TextCampaign

	for index := range body.TargetNumbers {
		params := &api.CreateMessageParams{}
		params.SetBody(body.Body)
		params.SetFrom(organization.TwilioNumber)
		params.SetTo(numbers[index])

		resp, err := client.Api.CreateMessage(params)

		if err != nil || resp.ErrorCode != nil {
			log.Fatal(err)
		}
	}
}
