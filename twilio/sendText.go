package twilio

import (
	"log"
	"os"
	"time"

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

	var numbers []string

	// get the numbers associated with the customer ids in the to field of the body

	var body Body
	err = c.BindJSON(&body)
	if err != nil {
		return
	}

	for _, number := range body.To {
		var customer model.Customer
		config.DB.Where("id = ?", number).First(&customer)
		if customer.PhoneNumber == "" {
			return
		}
		numbers = append(numbers, customer.PhoneNumber)
		customer.DatesReceivedSMS = append(customer.DatesReceivedSMS, time.Now().Format("2006-01-02 15:04:05"))
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
			"message":  "success",
			"response": resp,
		})

	}
}

func SendScheduledTexts(TextCampaign model.TextCampaign, org model.Organization) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	err = os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	if err != nil {
		return
	}
	err = os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))
	if err != nil {
		return
	}

	client := twilio.NewRestClient()

	organization := org

	if organization.TwilioNumber == "" {
		return
	}

	numbers := TextCampaign.TargetNumbers
	body := TextCampaign

	for index := range body.TargetNumbers {
		params := &api.CreateMessageParams{}
		params.SetBody(body.Body)
		params.SetFrom(organization.TwilioNumber)
		params.SetTo(numbers[index])

		resp, err := client.Api.CreateMessage(params)

		if err != nil || resp.ErrorCode != nil {
			return
		}
	}
}

func SendFlowTexts(ids []int32, org model.Organization, textBody string) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file for twilio")
	}

	err = os.Setenv("TWILIO_ACCOUNT_SID", os.Getenv("TWILIO_ACCOUNT_SID"))
	if err != nil {
		return
	}
	err = os.Setenv("TWILIO_AUTH_TOKEN", os.Getenv("TWILIO_AUTH_TOKEN"))
	if err != nil {
		return
	}

	client := twilio.NewRestClient()

	organization := org

	if organization.TwilioNumber == "" {
		return
	}

	var numbers []string
	var customers []model.Customer
	config.DB.Where("id IN ?", ids).Find(&customers)

	timeNowString := time.Now().Format("2006-01-02 15:04:05")

	for _, customer := range customers {
		numbers = append(numbers, customer.PhoneNumber)
		var customerTarget model.Customer
		config.DB.Where("id = ?", customer.ID).First(&customerTarget)
		customerTarget.DatesReceivedSMS = append(customerTarget.DatesReceivedSMS, timeNowString)
		config.DB.Save(&customerTarget)
	}

	for index := range numbers {
		params := &api.CreateMessageParams{}
		params.SetBody(textBody)
		params.SetFrom(organization.TwilioNumber)
		params.SetTo(numbers[index])

		resp, err := client.Api.CreateMessage(params)

		if err != nil || resp.ErrorCode != nil {
			return
		}
	}
}
